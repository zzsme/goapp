package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gowk/internal/app/errors"
)

var (
	// Logger instances for different log levels
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
)

// InitLogger initializes the logging system
func InitLogger() error {
	// Create logs directory if it doesn't exist
	logDir := filepath.Dir(ConfigData.Log.File)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return errors.Wrap(errors.Config, "failed to create log directory", err)
	}

	// Open log file
	file, err := os.OpenFile(ConfigData.Log.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(errors.Config, "failed to open log file", err)
	}

	// Initialize loggers with different prefixes
	infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger = log.New(file, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

// Info logs an info message with optional key-value pairs
func Info(msg string, keyvals ...interface{}) {
	logMessage(infoLogger, msg, keyvals...)
}

// Warn logs a warning message with optional key-value pairs
func Warn(msg string, keyvals ...interface{}) {
	logMessage(warnLogger, msg, keyvals...)
}

// Error logs an error message with optional key-value pairs
func Error(msg string, keyvals ...interface{}) {
	logMessage(errorLogger, msg, keyvals...)
}

// logMessage formats and writes a log message with key-value pairs
func logMessage(logger *log.Logger, msg string, keyvals ...interface{}) {
	// Format key-value pairs
	var details string
	if len(keyvals) > 0 {
		details = " |"
		for i := 0; i < len(keyvals); i += 2 {
			key := keyvals[i]
			var value interface{} = "MISSING"
			if i+1 < len(keyvals) {
				value = keyvals[i+1]
			}
			details += fmt.Sprintf(" %v=%v", key, value)
		}
	}

	logger.Printf("%s%s", msg, details)
}

// Now returns the current time
func Now() time.Time {
	return time.Now()
}
