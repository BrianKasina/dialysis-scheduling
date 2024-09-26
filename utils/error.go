package utils

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
)

// ErrorHandler function for handling errors and exceptions
func ErrorHandler(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    statusCode,
		"message": err.Error(),
		"stack":   string(debug.Stack()), // Include stack trace for debugging
	})
}
