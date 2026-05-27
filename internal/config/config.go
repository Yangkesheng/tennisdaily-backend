package config

import (
	"fmt"
	"os"
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
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslMode"`
	TimeZone string `yaml:"timeZone"`
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
	path := "config.yaml"
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

func (c databaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		c.Host,
		c.User,
		c.Password,
		c.Name,
		c.Port,
		c.SSLMode,
		c.TimeZone,
	)
}
