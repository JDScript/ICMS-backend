package authentication

import (
	"icms/internal/component/facial"
	"icms/internal/component/jwt"
	"icms/internal/repository/user"
)

type AuthenticationHandler struct {
	userRepo *user.UserRepository
	jwt      *jwt.Jwt
	facial   *facial.Facial
}

func NewHandler(userRepo *user.UserRepository, jwt *jwt.Jwt, facial *facial.Facial) *AuthenticationHandler {
	return &AuthenticationHandler{
		userRepo: userRepo,
		jwt:      jwt,
		facial:   facial,
	}
}
