package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"tennisdaily-backend/internal/logger"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port                    string
	DatabaseDSN             string
	JWTSecret               string
	JWTExpire               time.Duration
	WechatAppID             string
	WechatSecret            string
	StorageEnabled          bool
	StorageEnvID            string
	StorageBucket           string
	StorageRegion           string
	StorageFolder           string
	StorageSchedule         string
	StorageRunOnStart       bool
	SessionCategoryResolver *SessionCategoryResolver
	PolyesterStringHealth   *PolyesterStringHealthResolver
}

type fileConfig struct {
	Server                serverConfig                `yaml:"server"`
	Database              databaseConfig              `yaml:"database"`
	JWT                   jwtConfig                   `yaml:"jwt"`
	Wechat                wechatConfig                `yaml:"wechat"`
	Logger                loggerConfig                `yaml:"logger"`
	Storage               storageConfig               `yaml:"storage"`
	SessionEnums          SessionEnumsConfig          `yaml:"sessionEnums"`
	PolyesterStringHealth PolyesterStringHealthConfig `yaml:"polyesterStringHealth"`
}

type serverConfig struct {
	Port string `yaml:"port"`
}

type databaseConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Name      string `yaml:"name"`
	Charset   string `yaml:"charset"`
	ParseTime bool   `yaml:"parseTime"`
	Loc       string `yaml:"loc"`
	TimeZone  string `yaml:"timeZone"`
}

type jwtConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expireHours"`
}

type wechatConfig struct {
	AppID     string `yaml:"appId"`
	AppSecret string `yaml:"appSecret"`
}

type storageConfig struct {
	Enabled    bool   `yaml:"enabled"`
	EnvID      string `yaml:"envId"`
	Bucket     string `yaml:"bucket"`
	Region     string `yaml:"region"`
	Folder     string `yaml:"folder"`
	Schedule   string `yaml:"schedule"`
	RunOnStart bool   `yaml:"runOnStart"`
}

type loggerConfig struct {
	Level string `yaml:"level"`
}

func Load() Config {
	fc := loadConfigFile()

	logLevel := fc.Logger.Level
	logLevelSource := "config"
	if logLevel == "" {
		logLevel = "debug"
		logLevelSource = "default"
	}
	if value := os.Getenv("LOG_LEVEL"); value != "" {
		logLevel = value
		logLevelSource = "env"
	}
	if err := logger.SetLevel(logLevel); err != nil {
		panic(fmt.Errorf("invalid LOG_LEVEL: %s", logLevel))
	}

	port, portSource := envOrDefaultWithSource("SERVER_PORT", fc.Server.Port)
	databaseDSN, databaseDSNSource := envOrDefaultWithSource("DATABASE_DSN", fc.Database.DSN())
	jwtSecret, jwtSecretSource := envOrDefaultWithSource("JWT_SECRET", fc.JWT.Secret)
	wechatAppID, wechatAppIDSource := envOrDefaultWithSource("WECHAT_APP_ID", fc.Wechat.AppID)
	wechatSecret, wechatSecretSource := envOrDefaultWithSource("WECHAT_APP_SECRET", fc.Wechat.AppSecret)
	jwtExpireHours := fc.JWT.ExpireHours
	jwtExpireHoursSource := "config"
	if value := os.Getenv("JWT_EXPIRE_HOURS"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed <= 0 {
			panic(fmt.Errorf("invalid JWT_EXPIRE_HOURS: %s", value))
		}
		jwtExpireHours = parsed
		jwtExpireHoursSource = "env"
	}

	storageEnabled := fc.Storage.Enabled
	storageEnabledSource := "config"
	if value := os.Getenv("STORAGE_ENABLED"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			panic(fmt.Errorf("invalid STORAGE_ENABLED: %s", value))
		}
		storageEnabled = parsed
		storageEnabledSource = "env"
	}

	storageEnvID, _ := envOrDefaultWithSource("STORAGE_ENV_ID", fc.Storage.EnvID)
	storageBucket, _ := envOrDefaultWithSource("STORAGE_BUCKET", fc.Storage.Bucket)
	storageRegion, _ := envOrDefaultWithSource("STORAGE_REGION", fc.Storage.Region)
	storageFolder, _ := envOrDefaultWithSource("STORAGE_FOLDER", fc.Storage.Folder)
	storageSchedule, _ := envOrDefaultWithSource("STORAGE_SCHEDULE", fc.Storage.Schedule)

	storageRunOnStart := fc.Storage.RunOnStart
	if value := os.Getenv("STORAGE_RUN_ON_START"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			panic(fmt.Errorf("invalid STORAGE_RUN_ON_START: %s", value))
		}
		storageRunOnStart = parsed
	}

	sessionCategoryResolver, err := NewSessionCategoryResolver(fc.SessionEnums)
	if err != nil {
		panic(fmt.Errorf("invalid sessionEnums config: %w", err))
	}

	polyesterStringHealth, err := NewPolyesterStringHealthResolver(fc.PolyesterStringHealth)
	if err != nil {
		panic(fmt.Errorf("invalid polyesterStringHealth config: %w", err))
	}

	logger.Debug("config loaded logger.level source=%s value=%s", logLevelSource, logLevel)
	logger.Debug("config loaded server.port source=%s value=%s", portSource, port)
	logger.Debug("config loaded database.dsn source=%s value=%s", databaseDSNSource, maskDSN(databaseDSN))
	logger.Debug("config loaded jwt.secret source=%s value=%s", jwtSecretSource, maskSecret(jwtSecret))
	logger.Debug("config loaded jwt.expireHours source=%s value=%d", jwtExpireHoursSource, jwtExpireHours)
	logger.Debug("config loaded wechat.appId source=%s value=%s", wechatAppIDSource, maskMiddle(wechatAppID))
	logger.Debug("config loaded wechat.appSecret source=%s value=%s", wechatSecretSource, maskSecret(wechatSecret))
	logger.Debug("config loaded storage.enabled source=%s value=%t", storageEnabledSource, storageEnabled)
	logger.Debug("config loaded storage.envId source=%s value=%s", "config", maskMiddle(storageEnvID))
	logger.Debug("config loaded storage.bucket source=%s value=%s", "config", maskMiddle(storageBucket))
	logger.Debug("config loaded storage.region source=%s value=%s", "config", storageRegion)
	logger.Debug("config loaded storage.folder source=%s value=%s", "config", storageFolder)
	logger.Debug("config loaded storage.schedule source=%s value=%s", "config", storageSchedule)
	logger.Debug("config loaded storage.runOnStart source=%s value=%t", "config", storageRunOnStart)

	return Config{
		Port:                    port,
		DatabaseDSN:             databaseDSN,
		JWTSecret:               jwtSecret,
		JWTExpire:               time.Duration(jwtExpireHours) * time.Hour,
		WechatAppID:             wechatAppID,
		WechatSecret:            wechatSecret,
		StorageEnabled:          storageEnabled,
		StorageEnvID:            storageEnvID,
		StorageBucket:           storageBucket,
		StorageRegion:           storageRegion,
		StorageFolder:           storageFolder,
		StorageSchedule:         storageSchedule,
		StorageRunOnStart:       storageRunOnStart,
		SessionCategoryResolver: sessionCategoryResolver,
		PolyesterStringHealth:   polyesterStringHealth,
	}
}

func envOrDefaultWithSource(key string, fallback string) (string, string) {
	if value := os.Getenv(key); value != "" {
		return value, "env"
	}
	return fallback, "config"
}

func maskSecret(value string) string {
	if value == "" {
		return "<empty>"
	}
	return fmt.Sprintf("<set length=%d>", len(value))
}

func maskMiddle(value string) string {
	if value == "" {
		return "<empty>"
	}
	if len(value) <= 8 {
		return fmt.Sprintf("%s***", value[:1])
	}
	return fmt.Sprintf("%s***%s", value[:4], value[len(value)-4:])
}

func maskDSN(value string) string {
	if value == "" {
		return "<empty>"
	}

	atIndex := strings.Index(value, "@")
	if atIndex <= 0 {
		return maskMiddle(value)
	}

	return fmt.Sprintf("<credentials-hidden>%s", value[atIndex:])
}

func loadConfigFile() fileConfig {
	path := configPath()
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("read config file %s: %w", path, err))
	}

	var fc fileConfig
	if err := yaml.Unmarshal(content, &fc); err != nil {
		panic(fmt.Errorf("parse config file %s: %w", path, err))
	}
	return fc
}

func configPath() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}

	if _, err := os.Stat("config.yaml"); err == nil {
		return "config.yaml"
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "config.yaml"
	}

	return filepath.Join(filepath.Dir(filename), "..", "..", "config.yaml")
}

func (c databaseConfig) DSN() string {
	charset := c.Charset
	if charset == "" {
		charset = "utf8mb4"
	}

	loc := c.Loc
	if loc == "" {
		loc = c.TimeZone
	}
	if loc == "" {
		loc = "Asia/Shanghai"
	}

	parseTime := "False"
	if c.ParseTime || c.Loc == "" {
		parseTime = "True"
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		charset,
		parseTime,
		url.QueryEscape(loc),
	)
}
