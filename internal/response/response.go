package response

import (
	"encoding/json"
	"net/http"
)

type Body struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// Created
func Created(w http.ResponseWriter, msg string) {
	b := Body{
		Status:  http.StatusCreated,
		Message: msg,
	}

	sendResponse(w, &b)
}

// Success
func Success(w http.ResponseWriter, data interface{}, msg string) {
	b := Body{
		Status:  http.StatusOK,
		Data:    data,
		Message: msg,
	}

	sendResponse(w, &b)
}

// BadRequest
func BadRequest(w http.ResponseWriter, msg string, errs interface{}) {
	b := Body{
		Status:  http.StatusBadRequest,
		Message: msg,
		Errors:  errs,
	}

	sendResponse(w, &b)
}

func InternalServerError(w http.ResponseWriter, msg string) {
	b := Body{
		Status:  http.StatusInternalServerError,
		Message: msg,
	}

	sendResponse(w, &b)
}

// sendResponse
func sendResponse(w http.ResponseWriter, b *Body) {
	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(b.Status)

	json.NewEncoder(w).Encode(b)
}

// NotFound
func NotFound(w http.ResponseWriter, msg string) {
	b := Body{
		Status:  http.StatusNotFound,
		Message: msg,
	}

	sendResponse(w, &b)
}

// NotAllowded
func NotAllowded(w http.ResponseWriter, msg string) {
	b := Body{
		Status:  http.StatusMethodNotAllowed,
		Message: msg,
	}

	sendResponse(w, &b)
}
