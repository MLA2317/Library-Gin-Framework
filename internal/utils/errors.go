package utils

import (
	"fmt"
	"net/http"
)

// AppError represents a custom application error
type AppError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Error codes
const (
	ErrCodeValidation      = 1001
	ErrCodeNotFound        = 1002
	ErrCodeUnauthorized    = 1003
	ErrCodeForbidden       = 1004
	ErrCodeAlreadyExists   = 1005
	ErrCodeInternalServer  = 1006
	ErrCodeBadRequest      = 1007
	ErrCodeInvalidInput    = 1008
)

// Common error constructors
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
	}
}

func NewUnauthorizedError(message string) *AppError {
	if message == "" {
		message = "unauthorized"
	}
	return &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string) *AppError {
	if message == "" {
		message = "forbidden: insufficient permissions"
	}
	return &AppError{
		Code:       ErrCodeForbidden,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func NewAlreadyExistsError(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeAlreadyExists,
		Message:    fmt.Sprintf("%s already exists", resource),
		StatusCode: http.StatusConflict,
	}
}

func NewInternalServerError(message string, err error) *AppError {
	if message == "" {
		message = "internal server error"
	}
	return &AppError{
		Code:       ErrCodeInternalServer,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeBadRequest,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// WrapError wraps an existing error with additional context
func WrapError(err error, message string, statusCode int) *AppError {
	return &AppError{
		Code:       ErrCodeInternalServer,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}
