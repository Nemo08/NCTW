package router

import (
	"encoding/json"
	"net/http"
)

//Message хз зачем
func Message(message string) map[string]interface{} {
	return map[string]interface{}{"message": message}
}

//Respond формирует ответ HTTP
func Respond(w http.ResponseWriter, data map[string]interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
