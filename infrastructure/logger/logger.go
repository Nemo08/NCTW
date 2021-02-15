package logger

import (
	"go.uber.org/zap"
)

type LogInterface interface {
	WithField(key string, value interface{}) *Logger
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}
type Logger struct {
	z zap.Logger
}

func NewLogger() *Logger {
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	zl, _ := zap.NewProduction()
	defer zl.Sync()
	return &Logger{
		z: *zl,
	}
}

func (lg *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		z: *lg.z.With(zap.String(key, value.(string))),
	}
}

func (lg *Logger) Info(msg string) {
	lg.z.Info(msg)
}

func (lg *Logger) Warn(msg string) {
	lg.z.Warn(msg)
}

func (lg *Logger) Error(msg string) {
	lg.z.Error(msg)
}
