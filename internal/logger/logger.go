// Package logger предоставляет структурированное логирование.
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// Level represents log level.
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging.
type Logger struct {
	mu       sync.RWMutex
	out      io.Writer
	minLevel Level
	module   string
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// Init initializes the default logger.
func Init(out io.Writer, minLevel Level) {
	once.Do(func() {
		defaultLogger = &Logger{
			out:      out,
			minLevel: minLevel,
		}
		log.SetOutput(out)
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	})
}

// Default returns the default logger.
func Default() *Logger {
	if defaultLogger == nil {
		Init(os.Stdout, INFO)
	}
	return defaultLogger
}

// With creates a new logger with a module name.
func (l *Logger) With(module string) *Logger {
	return &Logger{
		out:      l.out,
		minLevel: l.minLevel,
		module:   module,
	}
}

// log logs a message at the given level.
func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.minLevel {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	module := ""
	if l.module != "" {
		module = fmt.Sprintf("[%s] ", l.module)
	}

	msg := fmt.Sprintf(format, args...)
	_ = fmt.Sprintf("%s %s%-5s %s\n", timestamp, module, level, msg) // Format for consistency

	log.Print(msg)
}

// Debug logs a debug message.
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message.
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message.
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs an error message.
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits.
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
	os.Exit(1)
}

// WithField logs with a key-value pair.
func (l *Logger) WithField(key string, value interface{}) *FieldLogger {
	return &FieldLogger{
		logger: l,
		key:    key,
		value:  value,
	}
}

// FieldLogger adds structured fields to logs.
type FieldLogger struct {
	logger *Logger
	key    string
	value  interface{}
}

func (f *FieldLogger) log(level Level, format string, args ...interface{}) {
	prefix := fmt.Sprintf("%s=%v ", f.key, f.value)
	msg := prefix + fmt.Sprintf(format, args...)
	f.logger.log(level, msg)
}

func (f *FieldLogger) Debug(format string, args ...interface{}) {
	f.log(DEBUG, format, args...)
}

func (f *FieldLogger) Info(format string, args ...interface{}) {
	f.log(INFO, format, args...)
}

func (f *FieldLogger) Warn(format string, args ...interface{}) {
	f.log(WARN, format, args...)
}

func (f *FieldLogger) Error(format string, args ...interface{}) {
	f.log(ERROR, format, args...)
}
