package types

type SendOTPRequest struct {
	Phone string `json:"phone"`
}

type VerifyOTPRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type UpdateProfileRequest struct {
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar_url"`
}

type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type UserResponse struct {
	ID         int64  `json:"id"`
	Phone      string `json:"phone"`
	FullName   string `json:"full_name"`
	AvatarUrl  string `json:"avatar_url"`
	IsVerified bool   `json:"is_verified"`
}
