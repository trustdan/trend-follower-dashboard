// Package responses provides standard HTTP response helpers
package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

// SuccessResponse wraps successful API responses
type SuccessResponse struct {
	Data interface{} `json:"data"`
}

// ErrorResponse wraps error API responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// JSON writes a JSON response with the given status code
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// Success writes a successful JSON response (200 OK)
func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, SuccessResponse{Data: data})
}

// Created writes a created response (201 Created)
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, SuccessResponse{Data: data})
}

// NoContent writes a no content response (204 No Content)
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Error writes an error JSON response
func Error(w http.ResponseWriter, status int, err error) {
	response := ErrorResponse{
		Error: http.StatusText(status),
		Code:  status,
	}

	if err != nil {
		response.Message = err.Error()
	}

	JSON(w, status, response)
}

// BadRequest writes a 400 Bad Request error
func BadRequest(w http.ResponseWriter, err error) {
	Error(w, http.StatusBadRequest, err)
}

// NotFound writes a 404 Not Found error
func NotFound(w http.ResponseWriter, err error) {
	Error(w, http.StatusNotFound, err)
}

// InternalError writes a 500 Internal Server Error
func InternalError(w http.ResponseWriter, err error) {
	Error(w, http.StatusInternalServerError, err)
}
