package logger

import (
	log "github.com/sirupsen/logrus"
)

type stdLog struct {
}

func NewStdLogger() stdLog {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	return stdLog{}
}

func (sl stdLog) LogMessage(v ...interface{}) {
	log.Infoln(v)
}

func (sl stdLog) LogError(v ...interface{}) {
	log.Errorln(v)
}

func (sl stdLog) Print(v ...interface{}) {
	log.Infoln(v)
}

func (sl stdLog) Write(b []byte) (int, error) {
	log.Infoln("[Сервер статики: ", string(b), "]")
	return len(b), nil
}
