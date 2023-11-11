package me

import (
	"errors"
	"icms/internal/model"
	"icms/internal/transport/http/request"
	"icms/pkg/auth"
	"icms/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func (handler *MeHandler) GetEnrolments(c *gin.Context) {
	userId := auth.CurrentUID(c)
	enrolments := handler.enrolmentRepo.EnrolledCourses(userId)

	response.JSON(c, http.StatusOK, true, "Query success", enrolments)
}

func (handler *MeHandler) CreateEnrolment(c *gin.Context) {
	req := request.MeEnrolCourseRequest{}
	if ok := request.BindAndValidate(c, &req); !ok {
		return
	}

	userId := auth.CurrentUID(c)
	enrolment := model.Enrolment{
		CourseID: req.CourseID,
		UserID:   userId,
	}

	err := handler.enrolmentRepo.Create(&enrolment)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			response.BadRequest(c, mysqlErr, "You already enrolled in this course")
			return
		}

		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1452 {
			response.BadRequest(c, mysqlErr, "Specified course not exists")
			return
		}

		response.Abort500(c, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, true, "Enroled successfully", enrolment)
}
