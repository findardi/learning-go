package common

import (
	"errors"
	"net/http"
)

// Application Errors
var (
	// Authentication & Authorization Errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrTokenExpired       = errors.New("token has expired")
	ErrTokenInvalid       = errors.New("invalid token")
	ErrTokenMissing       = errors.New("token is missing")
	ErrTokenGenerate      = errors.New("token generated is failed")

	// Duplicate/Conflict Errors
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUsername = errors.New("username already exists")
	ErrConflict          = errors.New("resource conflict")

	// Not Found Errors
	ErrNotFound     = errors.New("resource not found")
	ErrUserNotFound = errors.New("user not found")

	// Validation Errors
	ErrValidation   = errors.New("validation failed")
	ErrInvalidJSON  = errors.New("invalid JSON format")
	ErrInvalidInput = errors.New("invalid input")
	ErrBadRequest   = errors.New("bad request")

	// Password Errors
	ErrGeneratePassword = errors.New("failed to generate hash password")
	ErrPasswordMismatch = errors.New("password does not match")
	ErrWeakPassword     = errors.New("password is too weak")

	// Database Errors
	ErrDatabase       = errors.New("database error")
	ErrNoRowsAffected = errors.New("no rows affected")
	ErrTransaction    = errors.New("transaction failed")

	// Rate Limiting & Throttling
	ErrTooManyRequests   = errors.New("too many requests")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	// Internal Server Errors
	ErrInternalServer     = errors.New("internal server error")
	ErrServiceUnavailable = errors.New("service unavailable")

	// File & Upload Errors
	ErrFileTooBig      = errors.New("file size exceeds limit")
	ErrInvalidFileType = errors.New("invalid file type")
	ErrFileUpload      = errors.New("file upload failed")
)

// HTTP Status Code to Error Message mapping
var GetErrorMessages = map[int]string{
	http.StatusBadRequest:          "BAD_REQUEST",
	http.StatusUnauthorized:        "UNAUTHORIZED",
	http.StatusForbidden:           "FORBIDDEN",
	http.StatusNotFound:            "NOT_FOUND",
	http.StatusConflict:            "CONFLICT",
	http.StatusUnprocessableEntity: "UNPROCESSABLE_ENTITY",
	http.StatusTooManyRequests:     "TOO_MANY_REQUESTS",
	http.StatusInternalServerError: "INTERNAL_SERVER_ERROR",
	http.StatusServiceUnavailable:  "SERVICE_UNAVAILABLE",
}

func GetErrorStatusCode(err error) int {
	switch {
	// 400 Bad Request
	case errors.Is(err, ErrInvalidJSON),
		errors.Is(err, ErrValidation),
		errors.Is(err, ErrInvalidInput),
		errors.Is(err, ErrBadRequest),
		errors.Is(err, ErrGeneratePassword),
		errors.Is(err, ErrTokenGenerate),
		errors.Is(err, ErrWeakPassword):
		return http.StatusBadRequest

	// 401 Unauthorized
	case errors.Is(err, ErrInvalidCredentials),
		errors.Is(err, ErrUnauthorized),
		errors.Is(err, ErrTokenExpired),
		errors.Is(err, ErrTokenInvalid),
		errors.Is(err, ErrTokenMissing):
		return http.StatusUnauthorized

	// 403 Forbidden
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden

	// 404 Not Found
	case errors.Is(err, ErrNotFound),
		errors.Is(err, ErrUserNotFound):
		return http.StatusNotFound

	// 409 Conflict
	case errors.Is(err, ErrDuplicateEmail),
		errors.Is(err, ErrDuplicateUsername),
		errors.Is(err, ErrConflict):
		return http.StatusConflict

	// 413 Payload Too Large
	case errors.Is(err, ErrFileTooBig):
		return http.StatusRequestEntityTooLarge

	// 422 Unprocessable Entity
	case errors.Is(err, ErrPasswordMismatch),
		errors.Is(err, ErrInvalidFileType):
		return http.StatusUnprocessableEntity

	// 429 Too Many Requests
	case errors.Is(err, ErrTooManyRequests),
		errors.Is(err, ErrRateLimitExceeded):
		return http.StatusTooManyRequests

	// 503 Service Unavailable
	case errors.Is(err, ErrServiceUnavailable):
		return http.StatusServiceUnavailable

	// 500 Internal Server Error (default)
	default:
		return http.StatusInternalServerError
	}
}
