package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// InitLogger initializes the global logger
func InitLogger() *logrus.Logger {
	Logger = logrus.New()

	// Set output to stdout
	Logger.SetOutput(os.Stdout)

	// Set log level
	Logger.SetLevel(logrus.InfoLevel)

	// Set JSON formatter for structured logging
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     false,
	})

	// You can also use text formatter for development
	// Logger.SetFormatter(&logrus.TextFormatter{
	// 	FullTimestamp:   true,
	// 	TimestampFormat: "2006-01-02 15:04:05",
	// })

	return Logger
}

// LogError logs an error with context
func LogError(err error, context string, fields map[string]interface{}) {
	if Logger == nil {
		InitLogger()
	}

	entry := Logger.WithFields(logrus.Fields(fields))
	entry.WithField("context", context).Error(err)
}

// LogInfo logs an informational message
func LogInfo(message string, fields map[string]interface{}) {
	if Logger == nil {
		InitLogger()
	}

	entry := Logger.WithFields(logrus.Fields(fields))
	entry.Info(message)
}

// LogWarning logs a warning message
func LogWarning(message string, fields map[string]interface{}) {
	if Logger == nil {
		InitLogger()
	}

	entry := Logger.WithFields(logrus.Fields(fields))
	entry.Warn(message)
}

// LogDebug logs a debug message
func LogDebug(message string, fields map[string]interface{}) {
	if Logger == nil {
		InitLogger()
	}

	entry := Logger.WithFields(logrus.Fields(fields))
	entry.Debug(message)
}
