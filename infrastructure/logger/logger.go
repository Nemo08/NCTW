package logger

import (
	logrus "github.com/sirupsen/logrus"
)

var Log Logger

type Logger struct {
	l logrus.Logger
}

func (lg *Logger) LogMessage(v ...interface{}) {
	lg.l.Infoln(v)
}

func (lg *Logger) LogError(v ...interface{}) {
	lg.l.Errorln(v)
}

func (lg *Logger) Print(v ...interface{}) {
	lg.l.Infoln(v)
}

func Write(b []byte) (int, error) {
	Log.Print("[Сервер статики: ", string(b), "]")
	return len(b), nil
}
