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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
