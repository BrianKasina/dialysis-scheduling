package utils

import (
	"encoding/json"
	"net/http"
	"runtime"
)

// ErrorHandler function for handling errors and exceptions
func ErrorHandler(w http.ResponseWriter, statusCode int, err error, message string) {
	//capture the file and line where the error occured
	_, file, line, _ := runtime.Caller(1)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":  statusCode,
		"error": err.Error(),
		"message": message,
		"file":   file,
		"line":   line,
	})
}
