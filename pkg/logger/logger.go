package logger

import (
	"fmt"
	"strings"
)

// LogLevel represents the logging level
type LogLevel int

const (
	ErrorLevel LogLevel = iota
	WarnLevel
	InfoLevel
	DebugLevel
)

var currentLevel = InfoLevel

// SetLevel sets the logging level from string
func SetLevel(level string) error {
	switch strings.ToLower(level) {
	case "error":
		currentLevel = ErrorLevel
	case "warn", "warning":
		currentLevel = WarnLevel
	case "info":
		currentLevel = InfoLevel
	case "debug":
		currentLevel = DebugLevel
	default:
		return fmt.Errorf("invalid log level: %s (valid: error, warn, info, debug)", level)
	}
	return nil
}

// Info logs an info message (shown by default)
func Info(args ...interface{}) {
	if currentLevel >= InfoLevel {
		fmt.Println(args...)
	}
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	if currentLevel >= InfoLevel {
		fmt.Printf(format+"\n", args...)
	}
}

// Debug logs a debug message (only shown with debug level)
func Debug(args ...interface{}) {
	if currentLevel >= DebugLevel {
		fmt.Print("[DEBUG] ")
		fmt.Println(args...)
	}
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	if currentLevel >= DebugLevel {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	if currentLevel >= WarnLevel {
		fmt.Print("[WARN] ")
		fmt.Println(args...)
	}
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	if currentLevel >= WarnLevel {
		fmt.Printf("[WARN] "+format+"\n", args...)
	}
}

// Error logs an error message (always shown)
func Error(args ...interface{}) {
	fmt.Print("[ERROR] ")
	fmt.Println(args...)
}

// Errorf logs a formatted error message (always shown)
func Errorf(format string, args ...interface{}) {
	fmt.Printf("[ERROR] "+format+"\n", args...)
}

// GetLevel returns the current log level as string
func GetLevel() string {
	switch currentLevel {
	case ErrorLevel:
		return "error"
	case WarnLevel:
		return "warn"
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	default:
		return "info"
	}
}
