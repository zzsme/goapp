package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"goapp/internal/context"

	"github.com/gin-gonic/gin"
)

var (
	logger *Logger
)

// Logger represents our custom logger
type Logger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

// InitLogger initializes the logger
func InitLogger() {
	logFile := ConfigData.Log.Filename

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(logFile), 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Open log file
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Create loggers for different levels
	logger = &Logger{
		debugLogger: log.New(file, "DEBUG: ", log.Ldate|log.Ltime),
		infoLogger:  log.New(file, "INFO:  ", log.Ldate|log.Ltime),
		warnLogger:  log.New(file, "WARN:  ", log.Ldate|log.Ltime),
		errorLogger: log.New(file, "ERROR: ", log.Ldate|log.Ltime),
	}
}

// GetTraceID retrieves the trace ID from the context
func GetTraceID(args ...interface{}) string {
	for i := 0; i < len(args)-1; i += 2 {
		if args[i] == "trace_id" {
			if traceID, ok := args[i+1].(string); ok {
				return traceID
			}
		}
	}
	return "-"
}

// formatMessage formats a log message with caller information and arguments
func formatMessage(msg string, args ...interface{}) string {
	// Get caller information
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Base(file)

	// Get trace ID
	traceID := GetTraceID(args...)

	// Format the message
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("[%s:%d][%s] ", file, line, traceID))
	builder.WriteString(msg)

	// Add key-value pairs if any
	if len(args) > 0 {
		builder.WriteString(" |")
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				// Skip trace_id as it's already included in the prefix
				if args[i] != "trace_id" {
					builder.WriteString(fmt.Sprintf(" %v=%v", args[i], args[i+1]))
				}
			}
		}
	}

	return builder.String()
}

// Debug logs a debug message
func Debug(msg string, args ...interface{}) {
	if ConfigData.Log.Level == "debug" {
		logger.debugLogger.Println(formatMessage(msg, args...))
	}
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	if ConfigData.Log.Level == "debug" {
		logger.debugLogger.Println(formatMessage(fmt.Sprintf(format, args...)))
	}
}

// Info logs an info message
func Info(msg string, args ...interface{}) {
	logger.infoLogger.Println(formatMessage(msg, args...))
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	logger.infoLogger.Println(formatMessage(fmt.Sprintf(format, args...)))
}

// Warn logs a warning message
func Warn(msg string, args ...interface{}) {
	logger.warnLogger.Println(formatMessage(msg, args...))
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	logger.warnLogger.Println(formatMessage(fmt.Sprintf(format, args...)))
}

// Error logs an error message
func Error(msg string, args ...interface{}) {
	logger.errorLogger.Println(formatMessage(msg, args...))
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	logger.errorLogger.Println(formatMessage(fmt.Sprintf(format, args...)))
}

// Log is an alias for Info for backward compatibility
func Log(msg string, args ...interface{}) {
	Info(msg, args...)
}

// Context-aware logging functions

// DebugContext logs a debug message with context
func DebugContext(ctx *gin.Context, msg string, args ...interface{}) {
	if ConfigData.Log.Level == "debug" {
		newArgs := appendRequestID(ctx, args...)
		logger.debugLogger.Println(formatMessage(msg, newArgs...))
	}
}

// InfoContext logs an info message with context
func InfoContext(ctx *gin.Context, msg string, args ...interface{}) {
	newArgs := appendRequestID(ctx, args...)
	logger.infoLogger.Println(formatMessage(msg, newArgs...))
}

// WarnContext logs a warning message with context
func WarnContext(ctx *gin.Context, msg string, args ...interface{}) {
	newArgs := appendRequestID(ctx, args...)
	logger.warnLogger.Println(formatMessage(msg, newArgs...))
}

// ErrorContext logs an error message with context
func ErrorContext(ctx *gin.Context, msg string, args ...interface{}) {
	newArgs := appendRequestID(ctx, args...)
	logger.errorLogger.Println(formatMessage(msg, newArgs...))
}

// appendRequestID adds the request ID from context to the args
func appendRequestID(ctx *gin.Context, args ...interface{}) []interface{} {
	requestID := context.GetRequestID(ctx)
	newArgs := make([]interface{}, 0, len(args)+2)
	newArgs = append(newArgs, "trace_id", requestID)
	newArgs = append(newArgs, args...)
	return newArgs
}
