package http

import (
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// APIError creates a consistent error response with path, code, and message.
// This wraps huma's error model to ensure all errors follow the same structure.
func APIError(status int, code, message string, errs ...error) huma.StatusError {
	detail := message
	if len(errs) > 0 && errs[0] != nil {
		detail = fmt.Sprintf("%s: %s", message, errs[0].Error())
	}
	return huma.NewError(status, detail, &huma.ErrorDetail{
		Message:  message,
		Location: code,
	})
}

// Common error constructors with consistent codes.

func ErrBadRequest(path, message string) huma.StatusError {
	return APIError(http.StatusBadRequest, path, message)
}

func ErrNotFound(path, message string) huma.StatusError {
	return APIError(http.StatusNotFound, path, message)
}

func ErrInternal(path, message string, err error) huma.StatusError {
	return APIError(http.StatusInternalServerError, path, message, err)
}

func ErrUnprocessable(path, message string, err error) huma.StatusError {
	return APIError(http.StatusUnprocessableEntity, path, message, err)
}
