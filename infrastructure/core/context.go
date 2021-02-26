package core

type ServiceContext struct {
	RequestData  interface{}
	ResponseData interface{}
	param        map[string]string
	Log *Logr
}
