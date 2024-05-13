package utils

import (
	"encoding/json"
	"net/http"
)

type PayloadResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func WriteJson(data any, message string, status int, w http.ResponseWriter) error {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(
		PayloadResponse{Message: message, Data: data},
	)
}

func ReturnOkStatus(data any, message string, status int, w http.ResponseWriter) {
	WriteJson(data, message, status, w)
}

func ReturnErrorStatus(err error, status int, w http.ResponseWriter) {
	WriteJson(
		map[string]string{"error": err.Error()},
		"Error!",
		status,
		w,
	)
}
