package course

import (
	"icms/internal/repository/course"
	"icms/internal/repository/enrolment"
	"icms/internal/repository/message"
)

type CourseHandler struct {
	courseRepo    *course.CourseRepository
	enrolmentRepo *enrolment.EnrolmentRepository
	messageRepo   *message.MessageRepository
}

func NewHandler(courseRepo *course.CourseRepository, enrolmentRepo *enrolment.EnrolmentRepository, messageRepo *message.MessageRepository) *CourseHandler {
	return &CourseHandler{
		courseRepo:    courseRepo,
		enrolmentRepo: enrolmentRepo,
		messageRepo:   messageRepo,
	}
}
