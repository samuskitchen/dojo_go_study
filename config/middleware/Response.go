package middleware

import (
	"encoding/json"
	"net/http"
)

// ErrorMessage standardized error response.
type ErrorMessage struct {
	Status  bool   `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Map is a convenient way to create objects of unknown types.
type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Map map[string]interface{}

// JSON standardized JSON response.
func JSON(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) error {
	if data == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statusCode)
		return nil
	}

	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write(j)
	return nil
}

// HTTPError standardized error response in JSON format.
func HTTPError(w http.ResponseWriter, r *http.Request, statusCode int, codeError string, message string) error {
	msg := ErrorMessage{
		Status:  false,
		Code:    codeError,
		Message: message,
	}

	return JSON(w, r, statusCode, msg)
}
