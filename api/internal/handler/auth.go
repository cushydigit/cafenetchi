package handler

import (
	"cafenetchi-api/internal/helpers"
	"cafenetchi-api/internal/logger"
	"cafenetchi-api/internal/service"
	"cafenetchi-api/internal/types"
	"net/http"
)

type Auth struct {
	svc    *service.Auth
	logger *logger.Logger
}

func NewAuth(svc *service.Auth, l *logger.Logger) *Auth {
	return &Auth{
		svc:    svc,
		logger: l,
	}
}

func (h *Auth) SendOTP(w http.ResponseWriter, r *http.Request) {
	var req types.SendOTPRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.Error(w, types.ErrInvalidRequest)
		return
	}

	// TODO: maybe require a middlewares for phone validation format
	if err := h.svc.SendOTP(r.Context(), req.Phone); err != nil {
		helpers.Error(w, types.ErrInternalServer)
		return
	}

	payload := types.Response{
		Error:   false,
		Message: "OPT send successfully",
	}

	helpers.OK(w, payload)

}

func (h *Auth) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req types.VerifyOTPRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.Error(w, types.ErrInvalidRequest)
		return
	}

	_, token, isNewUser, err := h.svc.ValidateOTP(r.Context(), req.Phone, req.Code)
	if err != nil {
		helpers.Error(w, types.ErrNotAuthenticated)
		return
	}

	payload := types.Response{
		Message: "Authentication Succeed!",
		Error:   false,
		Data:    token,
	}

	if isNewUser {
		payload.Message = "Account created successfully"
	} else {
		payload.Message = "Login successful"
	}

	helpers.OK(w, payload)

}
