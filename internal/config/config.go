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
	SessionCategoryResolver *SessionCategoryResolver
	PolyesterStringHealth   *PolyesterStringHealthResolver
}

type fileConfig struct {
	Server                serverConfig                `yaml:"server"`
	Database              databaseConfig              `yaml:"database"`
	JWT                   jwtConfig                   `yaml:"jwt"`
	Wechat                wechatConfig                `yaml:"wechat"`
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

func Load() Config {
	fc := loadConfigFile()

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

	sessionCategoryResolver, err := NewSessionCategoryResolver(fc.SessionEnums)
	if err != nil {
		panic(fmt.Errorf("invalid sessionEnums config: %w", err))
	}

	polyesterStringHealth, err := NewPolyesterStringHealthResolver(fc.PolyesterStringHealth)
	if err != nil {
		panic(fmt.Errorf("invalid polyesterStringHealth config: %w", err))
	}

	logger.Debug("config loaded server.port source=%s value=%s", portSource, port)
	logger.Debug("config loaded database.dsn source=%s value=%s", databaseDSNSource, maskDSN(databaseDSN))
	logger.Debug("config loaded jwt.secret source=%s value=%s", jwtSecretSource, maskSecret(jwtSecret))
	logger.Debug("config loaded jwt.expireHours source=%s value=%d", jwtExpireHoursSource, jwtExpireHours)
	logger.Debug("config loaded wechat.appId source=%s value=%s", wechatAppIDSource, maskMiddle(wechatAppID))
	logger.Debug("config loaded wechat.appSecret source=%s value=%s", wechatSecretSource, maskSecret(wechatSecret))

	return Config{
		Port:                    port,
		DatabaseDSN:             databaseDSN,
		JWTSecret:               jwtSecret,
		JWTExpire:               time.Duration(jwtExpireHours) * time.Hour,
		WechatAppID:             wechatAppID,
		WechatSecret:            wechatSecret,
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
