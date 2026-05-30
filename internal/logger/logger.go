package logger

import "log"

func Debug(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}
