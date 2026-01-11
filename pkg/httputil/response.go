package httputil

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/askme/api/internal/domain"
)

// Response is a standard API response wrapper with generic data type
type Response[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ErrorResponse is used for error responses without data
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// JSON writes a type-safe JSON response
func JSON[T any](w http.ResponseWriter, status int, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := Response[T]{
		Success: status >= 200 && status < 300,
		Data:    data,
	}

	json.NewEncoder(w).Encode(resp)
}

// Error writes an error response
func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ErrorResponse{
		Success: false,
		Error:   message,
	}

	json.NewEncoder(w).Encode(resp)
}

// ErrorFromDomain maps domain errors to HTTP status codes
func ErrorFromDomain(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrAlreadyExists):
		Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrInvalidInput):
		Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		Error(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domain.ErrForbidden):
		Error(w, http.StatusForbidden, err.Error())
	case errors.Is(err, domain.ErrMutualFollowRequired):
		Error(w, http.StatusForbidden, err.Error())
	default:
		Error(w, http.StatusInternalServerError, "internal server error")
	}
}
