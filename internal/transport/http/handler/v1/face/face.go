package face

import (
	"icms/internal/component/facial"
	"icms/internal/repository/user"
)

type FaceHandler struct {
	userRepo *user.UserRepository
	facial   *facial.Facial
}

func NewHandler(userRepo *user.UserRepository, facial *facial.Facial) *FaceHandler {
	return &FaceHandler{
		userRepo: userRepo,
		facial:   facial,
	}
}
