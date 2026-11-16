package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func Success(w http.ResponseWriter, data interface{}, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	JSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
		Message: msg,
	})
}

func Error(w http.ResponseWriter, statusCode int, errorMsg string) {
	JSON(w, statusCode, APIResponse{
		Success: false,
		Error:   errorMsg,
	})
}
