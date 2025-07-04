package middleware

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger initializes the logger with a simple configuration
func InitLogger() {

	// Adding simple config for development
	// Can be extended to support production
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create the logger
	var err error
	Logger, err = config.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	// Replace global logger
	zap.ReplaceGlobals(Logger)
}

// GetLogger returns the configured logger
func GetLogger() *zap.Logger {
	if Logger == nil {
		InitLogger()
	}
	return Logger
}

// Sync flushes any buffered log entries
func Sync() {
	if Logger != nil {
		Logger.Sync()
	}
} 