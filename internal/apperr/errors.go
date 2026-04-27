package apperr

import (
	"errors"
	"net/http"
)

var ErrNotFound = errors.New("resource not found")

type AppError struct {
	Code    int    // HTTP status code
	Message string // user-facing message
	Err     error  // underlying error for logging
}

func (e *AppError) Error() string { return e.Message }

func NotFound(msg string, err error) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg, Err: err}
}

func BadRequest(msg string, err error) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg, Err: err}
}

func Internal(err error) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: "Internal server error", Err: err}
}
