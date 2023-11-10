package user

import (
	"icms/internal/component/facial"
	"icms/internal/repository/user"
)

type UserHandler struct {
	userRepo *user.UserRepository
	facial   *facial.Facial
}

func NewHandler(userRepo *user.UserRepository, facial *facial.Facial) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
		facial:   facial,
	}
}
