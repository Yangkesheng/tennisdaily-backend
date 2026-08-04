package logger

import (
	"fmt"
	"log"
	"strings"
)

// Level 日志级别。
type Level int

const (
	LevelError Level = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

var levelNames = map[Level]string{
	LevelError: "ERROR",
	LevelWarn:  "WARN",
	LevelInfo:  "INFO",
	LevelDebug: "DEBUG",
}

// currentLevel 当前生效的日志级别，默认 debug，保持原有行为。
var currentLevel = LevelDebug

// SetLevel 设置全局日志级别，支持 debug/info/warn/error，不区分大小写。
// 低于该级别的日志将被过滤，例如 info 会过滤 debug。
func SetLevel(level string) error {
	normalized := strings.ToLower(strings.TrimSpace(level))
	switch normalized {
	case "debug":
		currentLevel = LevelDebug
	case "info":
		currentLevel = LevelInfo
	case "warn", "warning":
		currentLevel = LevelWarn
	case "error":
		currentLevel = LevelError
	default:
		return fmt.Errorf("invalid log level %q (expected debug/info/warn/error)", level)
	}
	return nil
}

// Enabled 判断指定级别当前是否会被输出。
func Enabled(level Level) bool {
	return level <= currentLevel
}

// Debug 输出 debug 级别日志。
func Debug(format string, args ...interface{}) {
	logf(LevelDebug, format, args...)
}

// Info 输出 info 级别日志。
func Info(format string, args ...interface{}) {
	logf(LevelInfo, format, args...)
}

// Warn 输出 warn 级别日志。
func Warn(format string, args ...interface{}) {
	logf(LevelWarn, format, args...)
}

// Error 输出 error 级别日志。
func Error(format string, args ...interface{}) {
	logf(LevelError, format, args...)
}

func logf(level Level, format string, args ...interface{}) {
	if !Enabled(level) {
		return
	}
	log.Printf("[%s] "+format, append([]interface{}{levelNames[level]}, args...)...)
}
