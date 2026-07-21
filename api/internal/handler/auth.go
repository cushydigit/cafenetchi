package handler

import (
	"cafenetchi-api/internal/helpers"
	"cafenetchi-api/internal/service"
	"cafenetchi-api/internal/types"
	"cafenetchi-api/internal/validation"
	"net/http"
)

type Auth struct {
	svc service.Auth
}

func NewAuth(s service.Auth) *Auth {
	return &Auth{
		svc: s,
	}
}

func (h *Auth) SendOTP(w http.ResponseWriter, r *http.Request) {
	var req types.SendOTPRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.Error(w, types.ErrInvalidRequest)
		return
	}

	phone, err := validation.ValidateSendOTP(req)
	if err != nil {
		helpers.Error(w, err)
		return
	}

	if err := h.svc.SendOTP(r.Context(), phone); err != nil {
		helpers.Error(w, err)
		return
	}

	payload := types.Response{
		Error:   false,
		Message: "OTP send successfully",
	}

	helpers.OK(w, payload)

}

func (h *Auth) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req types.VerifyOTPRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.Error(w, types.ErrInvalidRequest)
		return
	}

	phone, code, err := validation.ValidateVerifyOTP(req)
	if err != nil {
		helpers.Error(w, err)
		return
	}

	result, err := h.svc.ValidateOTP(r.Context(), phone, code)
	if err != nil {
		helpers.Error(w, err)
		return
	}

	status := http.StatusOK
	message := "Login successful"

	if result.IsNewUser {
		status = http.StatusCreated
		message = "Account created successfully"
	}

	// TODO: return access token
	_ = helpers.WriteJSON(w, status, types.Response{
		Message: message,
		Data:    result,
	})

}
