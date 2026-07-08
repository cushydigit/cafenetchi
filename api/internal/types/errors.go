package types

import "net/http"

type APIError struct {
	Status  int    `json:"status"`
	Code    string `json:"error"`
	Message string `json:"message"`
}

var (
	ErrInvalidRequest = APIError{
		Status:  http.StatusBadRequest,
		Code:    "INVALID_REQUEST",
		Message: "Invalid request",
	}

	ErrNotAuthenticated = APIError{
		Status:  http.StatusUnauthorized,
		Code:    "NOT_AUTHENTICATE",
		Message: "Not authenticated",
	}

	ErrNotFound = APIError{
		Status:  http.StatusNotFound,
		Code:    "NOT_FOUND",
		Message: "Not found",
	}

	ErrInvalidOTP = APIError{
		Status:  http.StatusBadRequest,
		Code:    "INVALID_OTP",
		Message: "Invalid OTP",
	}

	ErrInternalServer = APIError{
		Status:  http.StatusInternalServerError,
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Internal server error",
	}

	ErrTooManyRequest = APIError{
		Status:  http.StatusTooManyRequests,
		Code:    "RATE_LIMITED",
		Message: "too many requests",
	}
)
