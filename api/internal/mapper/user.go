package mapper

import (
	db "cafenetchi-api/internal/db/generated"
	"cafenetchi-api/internal/model"
	"cafenetchi-api/internal/types"
)

func DBUserToModel(u *db.User) *model.User {
	return &model.User{
		ID:        u.ID,
		FullName:  u.FullName.String,
		Phone:     u.Phone,
		AvatarURL: u.AvatarUrl.String,

		IsVerified: u.IsVerified.Bool,
		Status:     model.UserStatus(u.Status.String),
	}
}

func ModelUserToResponse(u *model.User) *types.UserResponse {
	return &types.UserResponse{
		ID:         u.ID,
		Phone:      u.Phone,
		FullName:   u.FullName,
		AvatarUrl:  u.AvatarURL,
		IsVerified: u.IsVerified,
	}
}
