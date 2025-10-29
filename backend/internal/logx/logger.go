package logx

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// Initialize sets up the global logger
func Initialize(logPath string) error {
	logger = logrus.New()

	// JSON formatting for structured logs
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	// Create log directory if it doesn't exist
	dir := filepath.Dir(logPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}
	}

	// Open log file
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Write logs only to file (not stdout)
	// Stdout is reserved for command output (JSON or human-readable)
	// Set TF_DEBUG=1 to also write logs to stderr for debugging
	if os.Getenv("TF_DEBUG") == "1" {
		logger.SetOutput(io.MultiWriter(os.Stderr, file))
	} else {
		logger.SetOutput(file)
	}

	// Set level from environment or default to INFO
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)

	return nil
}

// GenerateCorrelationID creates a new correlation ID
func GenerateCorrelationID() string {
	return uuid.New().String()
}

// WithCorrelationID creates a logger with correlation ID
func WithCorrelationID(corrID string) *logrus.Entry {
	if logger == nil {
		// Fallback to stderr if not initialized
		logger = logrus.New()
		logger.SetOutput(os.Stderr)
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	}
	return logger.WithField("corr_id", corrID)
}

// Info logs at info level
func Info(msg string, fields map[string]interface{}) {
	if logger == nil {
		logger = logrus.New()
		logger.SetOutput(os.Stderr)
	}
	logger.WithFields(fields).Info(msg)
}

// Error logs at error level
func Error(msg string, err error, fields map[string]interface{}) {
	if logger == nil {
		logger = logrus.New()
		logger.SetOutput(os.Stderr)
	}
	entry := logger.WithFields(fields)
	if err != nil {
		entry = entry.WithError(err)
	}
	entry.Error(msg)
}

// Debug logs at debug level
func Debug(msg string, fields map[string]interface{}) {
	if logger == nil {
		logger = logrus.New()
		logger.SetOutput(os.Stderr)
	}
	logger.WithFields(fields).Debug(msg)
}

// Warn logs at warn level
func Warn(msg string, fields map[string]interface{}) {
	if logger == nil {
		logger = logrus.New()
		logger.SetOutput(os.Stderr)
	}
	logger.WithFields(fields).Warn(msg)
}
