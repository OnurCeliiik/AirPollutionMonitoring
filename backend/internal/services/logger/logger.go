package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Level represents the severity level of a log entry
type Level int

const (
	// DEBUG level for detailed troubleshooting
	DEBUG Level = iota
	// INFO level for general operational information
	INFO
	// WARN level for events that might cause problems
	WARN
	// ERROR level for events that would impact functionality
	ERROR
	// FATAL level for events that would cause the service to stop
	FATAL
)

var levelNames = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

// Logger represents a structured logger
type Logger struct {
	level  Level
	prefix string
	out    io.Writer
}

// New creates a new logger with the given prefix and level
func New(prefix string, level Level) *Logger {
	return &Logger{
		level:  level,
		prefix: prefix,
		out:    os.Stdout,
	}
}

// WithOutput sets the output writer for the logger
func (l *Logger) WithOutput(out io.Writer) *Logger {
	l.out = out
	return l
}

// GetLevel returns the current logging level as a string
func (l *Logger) GetLevel() string {
	return levelNames[l.level]
}

// SetLevel sets the logging level from a string
func (l *Logger) SetLevel(levelStr string) {
	levelStr = strings.ToUpper(levelStr)
	for level, name := range levelNames {
		if name == levelStr {
			l.level = level
			return
		}
	}
	// Default to INFO if invalid level
	l.level = INFO
}

// log logs a message at the specified level
func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02T15:04:05.000Z")
	message := fmt.Sprintf(format, args...)
	logEntry := fmt.Sprintf("[%s] [%s] %s: %s\n", timestamp, levelNames[level], l.prefix, message)

	fmt.Fprint(l.out, logEntry)

	// If fatal, exit after logging
	if level == FATAL {
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

// DefaultLogger is the default logger instance
var DefaultLogger = New("APP", INFO)

// SetDefaultLogLevel sets the log level for the default logger
func SetDefaultLogLevel(level string) {
	DefaultLogger.SetLevel(level)
}

// Debug logs a debug message to the default logger
func Debug(format string, args ...interface{}) {
	DefaultLogger.Debug(format, args...)
}

// Info logs an info message to the default logger
func Info(format string, args ...interface{}) {
	DefaultLogger.Info(format, args...)
}

// Warn logs a warning message to the default logger
func Warn(format string, args ...interface{}) {
	DefaultLogger.Warn(format, args...)
}

// Error logs an error message to the default logger
func Error(format string, args ...interface{}) {
	DefaultLogger.Error(format, args...)
}

// Fatal logs a fatal message to the default logger and exits
func Fatal(format string, args ...interface{}) {
	DefaultLogger.Fatal(format, args...)
}
