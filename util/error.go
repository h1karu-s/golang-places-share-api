package util

import (
	"net/http"
	"encoding/json"
	"../models"
)

// Err handle error
func Err (w http.ResponseWriter, e error, statusCode int) {
	var error models.HTTPError
	error.Message = e.Error()
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(error)
}