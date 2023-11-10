package me

import (
	"icms/internal/component/facial"
	"icms/internal/repository/activity"
	"icms/internal/repository/enrolment"
	"icms/internal/repository/message"
	"icms/internal/repository/user"
	"icms/pkg/auth"
	"icms/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MeHandler struct {
	activityRepo  *activity.ActivityRepository
	enrolmentRepo *enrolment.EnrolmentRepository
	facial        *facial.Facial
	messageRepo   *message.MessageRepository
	userRepo      *user.UserRepository
}

func NewHandler(
	activityRepo *activity.ActivityRepository,
	enrolmentRepo *enrolment.EnrolmentRepository,
	facial *facial.Facial,
	messageRepo *message.MessageRepository,
	userRepo *user.UserRepository,
) *MeHandler {
	return &MeHandler{
		activityRepo:  activityRepo,
		enrolmentRepo: enrolmentRepo,
		facial:        facial,
		messageRepo:   messageRepo,
		userRepo:      userRepo,
	}
}

func (handler *MeHandler) GetMe(c *gin.Context) {
	user := auth.CurrentUser(c)
	response.JSON(c, http.StatusOK, true, "", user)
}
