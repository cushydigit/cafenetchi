package handler

import (
	"cafenetchi-api/internal/helpers"
	"cafenetchi-api/internal/service"
	"cafenetchi-api/internal/types"
	"errors"
	"log"
	"net/http"
)

type Auth struct {
	svc *service.Auth
}

func NewAuth(svc *service.Auth) *Auth {
	return &Auth{
		svc: svc,
	}
}

func (h *Auth) SendOTP(w http.ResponseWriter, r *http.Request) {
	var req types.SendOTPRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		log.Println(err)
		helpers.ErrorJSON(w, errors.New("Invalid request"))
		return
	}

	// TODO: maybe require a middlewares for phone validation format
	if err := h.svc.SendOTP(req.Phone); err != nil {
		log.Println(err)
		helpers.ErrorJSON(w, errors.New("Failed to send OTP"), http.StatusInternalServerError)
		return
	}

	payload := types.Response{
		Error:   false,
		Message: "OPT send successfully",
	}
	helpers.WriteJSON(w, http.StatusOK, payload)
}

func (h *Auth) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req types.VerifyOTPRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.ErrorJSON(w, errors.New("Invalid request"), http.StatusBadRequest)
		return
	}

	_, token, isNewUser, err := h.svc.ValidateOTP(req.Phone, req.Code)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("Not authenticated"), http.StatusUnauthorized)
		return
	}

	payload := types.Response{
		Message: "Authentication Succeed",
		Error:   false,
		Data:    token,
	}

	if isNewUser {
		payload.Message = "Account created successfully"
	} else {
		payload.Message = "Login successful"
	}

	helpers.WriteJSON(w, http.StatusOK, payload)

}
