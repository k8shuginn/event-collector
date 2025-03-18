package logger

import "go.uber.org/zap"

var (
	// global logger
	writer *zap.Logger
)

// Debug global logger debug level
func Debug(msg string, fields ...zap.Field) {
	writer.Debug(msg, fields...)
}

// Info global logger info level
func Info(msg string, fields ...zap.Field) {
	writer.Info(msg, fields...)
}

// Warn global logger warn level
func Warn(msg string, fields ...zap.Field) {
	writer.Warn(msg, fields...)
}

// Error global logger error level
func Error(msg string, fields ...zap.Field) {
	writer.Error(msg, fields...)
}

// Fatal global logger fatal level
func Fatal(msg string, fields ...zap.Field) {
	writer.Fatal(msg, fields...)
}

// Panic global logger panic level
func Panic(msg string, fields ...zap.Field) {
	writer.Panic(msg, fields...)
}
