package common

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var validate = validator.New()

func Validate(data any) []ValidationError {
	var errors []ValidationError

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var e ValidationError
			e.Field = strings.ToLower(err.Field())
			e.Message = getErrorMessage(err)
			errors = append(errors, e)
		}
	}

	return errors
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", err.Field(), err.Param())
	case "eqfield":
		return fmt.Sprintf("%s must be equal to %s", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s is not valid", err.Field())
	}
}

func ValidateRequest(data any) *ErrorResponse {
	if errors := Validate(data); len(errors) > 0 {
		return &ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors:  errors,
		}
	}
	return nil
}
