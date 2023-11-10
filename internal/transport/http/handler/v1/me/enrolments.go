package me

import (
	"icms/pkg/auth"
	"icms/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *MeHandler) GetEnrolments(c *gin.Context) {
	userId := auth.CurrentUID(c)
	enrolments := handler.enrolmentRepo.EnrolledCourses(userId)

	response.JSON(c, http.StatusOK, true, "Query success", enrolments)
}
