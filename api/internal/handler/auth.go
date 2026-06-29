package handler

import (
	"cafenetchi-api/internal/helpers"
	"cafenetchi-api/internal/service"
	"cafenetchi-api/internal/types"
	"errors"
	"log/slog"
	"net/http"
)

type Auth struct {
	svc    *service.Auth
	logger *slog.Logger
}

func NewAuth(svc *service.Auth, logger *slog.Logger) *Auth {
	return &Auth{
		svc:    svc,
		logger: logger,
	}
}

func (h *Auth) SendOTP(w http.ResponseWriter, r *http.Request) {
	var req types.SendOTPRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		h.errorJSON(w, errors.New("Invalid request"), http.StatusBadRequest)
		return
	}

	// TODO: maybe require a middlewares for phone validation format
	if err := h.svc.SendOTP(r.Context(), req.Phone); err != nil {
		h.errorJSON(w, errors.New("Failed to send OTP"), http.StatusInternalServerError)
		return
	}

	payload := types.Response{
		Error:   false,
		Message: "OPT send successfully",
	}

	h.writeJSON(w, http.StatusOK, payload)
}

func (h *Auth) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req types.VerifyOTPRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		h.errorJSON(w, errors.New("Invalid request"), http.StatusBadRequest)
		return
	}

	_, token, isNewUser, err := h.svc.ValidateOTP(r.Context(), req.Phone, req.Code)
	if err != nil {
		h.errorJSON(w, errors.New("Not authenticated"), http.StatusUnauthorized)
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

	h.writeJSON(w, http.StatusOK, payload)

}

func (h *Auth) errorJSON(w http.ResponseWriter, err error, status int) {
	if e := helpers.ErrorJSON(w, err, status); e != nil {
		h.logger.Error(
			"failed writing error json response",
			"error", e,
		)
	}
}

func (h *Auth) writeJSON(w http.ResponseWriter, status int, data any) {
	if e := helpers.WriteJSON(w, status, data); e != nil {
		h.logger.Error(
			"failed writing json response",
			"error", e,
		)
	}
}
