package usecase

type LogInterface interface {
	LogMessage(v ...interface{})
	LogError(v ...interface{})
	Print(v ...interface{})
	Write([]byte) (int, error)
}
