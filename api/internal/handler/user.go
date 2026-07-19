package handler

import (
	"cafenetchi-api/internal/helpers"
	"cafenetchi-api/internal/logger"
	"cafenetchi-api/internal/middleware"
	"cafenetchi-api/internal/model"
	"cafenetchi-api/internal/repository"
	"cafenetchi-api/internal/service"
	"cafenetchi-api/internal/types"
	"errors"
	"net/http"
)

type User struct {
	svc    service.User
	logger *logger.Logger
}

func NewUser(s service.User, l *logger.Logger) *User {
	return &User{
		svc:    s,
		logger: l,
	}
}

func (h *User) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	u, err := h.svc.GetByID(r.Context(), userID)
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		helpers.Error(w, types.ErrUserNotFound)
	default:
		helpers.Error(w, types.ErrInternalServer)
	}

	helpers.OK(w, u)
}

func (h *User) UpdateMe(w http.ResponseWriter, r *http.Request) {
	// TODO: add middleware for auth
	var req types.UpdateProfileRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.Error(w, types.ErrInvalidRequest)
		return
	}

	userID := middleware.UserID(r.Context())
	_, err := h.svc.UpdateProfile(r.Context(), userID, model.UpdateUser{
		FullName:  req.FullName,
		AvatarURL: req.Avatar,
	})
	if err != nil {
		helpers.Error(w, types.ErrInternalServer)
		return
	}

	helpers.OK(w, "user profile updated")
}
