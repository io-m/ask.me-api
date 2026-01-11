package httputil

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// DecodeJSON decodes JSON request body into the given struct
func DecodeJSON[T any](r *http.Request) (*T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

// PathValue extracts a path parameter from the request
func PathValue(r *http.Request, key string) string {
	return r.PathValue(key)
}

// QueryInt extracts an integer query parameter with a default value
func QueryInt(r *http.Request, key string, defaultVal int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

// QueryString extracts a string query parameter with a default value
func QueryString(r *http.Request, key string, defaultVal string) string {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}
	return val
}
