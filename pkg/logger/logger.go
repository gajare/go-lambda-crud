package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type Logger struct {
	Logrus *logrus.Logger
	Zap    *zap.Logger
}

func NewLogger() *Logger {
	// Logrus logger
	logrusLogger := logrus.New()
	logrusLogger.SetOutput(os.Stdout)
	logrusLogger.SetFormatter(&logrus.JSONFormatter{})

	// Zap logger
	zapLogger, _ := zap.NewProduction()

	return &Logger{
		Logrus: logrusLogger,
		Zap:    zapLogger,
	}
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logrus.Info(msg)
	l.Zap.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Logrus.Error(msg)
	l.Zap.Error(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logrus.Debug(msg)
	l.Zap.Debug(msg, fields...)
}
