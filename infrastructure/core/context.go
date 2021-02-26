package core

type ServiceContext struct {
	requestData  interface{}
	responseData interface{}
	param        map[string]string
	Log *Logr
}
