package logger

import (
	log "github.com/sirupsen/logrus"
)

func LogMessage(v ...interface{}) {
	log.Infoln(v)
}

func LogError(v ...interface{}) {
	log.Errorln(v)
}

func Print(v ...interface{}) {
	log.Infoln(v)
}

func Write(b []byte) (int, error) {
	log.Infoln("[Сервер статики: ", string(b), "]")
	return len(b), nil
}
