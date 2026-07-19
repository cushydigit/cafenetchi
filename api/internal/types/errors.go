package types

import "net/http"

type APIError struct {
	Status  int    `json:"status"`
	Code    string `json:"error"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
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

	ErrPhoneRequired = APIError{
		Status:  http.StatusBadRequest,
		Code:    "PHONE_REQUIRED",
		Message: "Phone number is required",
	}

	ErrOTPCodeRequired = APIError{
		Status:  http.StatusBadRequest,
		Code:    "CODE_REQUIRED",
		Message: "otp code required",
	}

	ErrInvalidPhone = APIError{
		Status:  http.StatusBadRequest,
		Code:    "INVALID_PHONE",
		Message: "Invalid iranian phone number",
	}

	ErrOTPRequired = APIError{
		Status:  http.StatusBadRequest,
		Code:    "OTP_REQUIRED",
		Message: "OTP is required",
	}

	ErrInvalidOTP = APIError{
		Status:  http.StatusBadRequest,
		Code:    "INVALID_OTP",
		Message: "Invalid OTP code",
	}

	ErrUserNotFound = APIError{
		Status:  http.StatusNotFound,
		Code:    "USER_NOT_FOUND",
		Message: "User not found",
	}

	ErrNameTooLong = APIError{
		Status:  http.StatusBadRequest,
		Code:    "NAME_TOO_LONG",
		Message: "Name is too long",
	}

	ErrNameTooShort = APIError{
		Status:  http.StatusBadRequest,
		Code:    "NAME_TOO_SHORT",
		Message: "Name is too short",
	}

	ErrAvatarURLTooLong = APIError{
		Status:  http.StatusBadRequest,
		Code:    "AVATAR_URL_TOO_LONG",
		Message: "Avatar URL is too long",
	}

	ErrAvatarURLTooShort = APIError{
		Status:  http.StatusBadRequest,
		Code:    "AVATAR_URL_TOO_SHORT",
		Message: "Avatar URL is too short",
	}

	ErrInvalidAvatarURL = APIError{
		Status:  http.StatusBadRequest,
		Code:    "INVALID_AVATAR_URL",
		Message: "Invalid avatar URL",
	}
)
