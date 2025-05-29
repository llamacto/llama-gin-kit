package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initializes the logger
func Init() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	var err error
	log, err = config.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
}

// Error logs an error message
func Error(msg string, err error) {
	if log == nil {
		Init()
	}
	log.Error(msg, zap.Error(err))
}

// Info logs an info message
func Info(msg string, args ...interface{}) {
	if log == nil {
		Init()
	}
	log.Sugar().Infof(msg, args...)
}

// Debug logs a debug message
func Debug(msg string, args ...interface{}) {
	if log == nil {
		Init()
	}
	log.Sugar().Debugf(msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...interface{}) {
	if log == nil {
		Init()
	}
	log.Sugar().Warnf(msg, args...)
}

// Fatal logs a fatal message and exits the program
func Fatal(msg string, args ...interface{}) {
	if log == nil {
		Init()
	}
	log.Sugar().Fatalf(msg, args...)
	os.Exit(1)
}

// Sync synchronizes the logger
func Sync() {
	if log != nil {
		log.Sync()
	}
}
