package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port         string
	DatabaseDSN  string
	JWTSecret    string
	JWTExpire    time.Duration
	WechatAppID  string
	WechatSecret string
}

type fileConfig struct {
	Server   serverConfig   `yaml:"server"`
	Database databaseConfig `yaml:"database"`
	JWT      jwtConfig      `yaml:"jwt"`
	Wechat   wechatConfig   `yaml:"wechat"`
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

	return Config{
		Port:         fc.Server.Port,
		DatabaseDSN:  fc.Database.DSN(),
		JWTSecret:    fc.JWT.Secret,
		JWTExpire:    time.Duration(fc.JWT.ExpireHours) * time.Hour,
		WechatAppID:  fc.Wechat.AppID,
		WechatSecret: fc.Wechat.AppSecret,
	}
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
