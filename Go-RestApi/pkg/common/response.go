package common

import (
	"encoding/json"
	"net/http"
	"time"
)

type SuccessResponse struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	Timestamp string `json:"timestamp"`
}

type ErrorResponse struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Errors    any    `json:"errors,omitempty"`
	Timestamp string `json:"timestamp"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	return nil
}

func SuccessResponseJSON(w http.ResponseWriter, status int, message string, data any) error {
	response := SuccessResponse{
		Status:    status,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return WriteJSON(w, status, response)
}

func ErrorResponseJSON(w http.ResponseWriter, status int, message string, errors any) error {
	response := ErrorResponse{
		Status:    status,
		Message:   message,
		Errors:    errors,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return WriteJSON(w, status, response)
}

func RespondWithSuccess(w http.ResponseWriter, message string, data any) error {
	return SuccessResponseJSON(w, http.StatusOK, message, data)
}

func RespondWithError(w http.ResponseWriter, err error) error {
	statusCode := GetErrorStatusCode(err)
	message := GetErrorMessages[statusCode]
	return ErrorResponseJSON(w, statusCode, message, err.Error())
}

func RespondWithCustomError(w http.ResponseWriter, status int, message string, errors any) error {
	return ErrorResponseJSON(w, status, message, errors)
}
