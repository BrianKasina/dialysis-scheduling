package utils

import (
	"encoding/json"
	"net/http"
	// "runtime/debug"
)

// ErrorHandler function for handling errors and exceptions
func ErrorHandler(w http.ResponseWriter, statusCode int, err error, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":  statusCode,
		"error": err.Error(),
		"message": message,
		// "stack":   string(debug.Stack()), // Include stack trace for debugging
	})
}
