package enrolment

import (
	"icms/internal/model"

	"gorm.io/gorm"
)

type EnrolmentRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *EnrolmentRepository {
	return &EnrolmentRepository{
		db: db,
	}
}

func (repo *EnrolmentRepository) EnrolledCourses(userId int32) []model.Enrolment {
	enrolments := make([]model.Enrolment, 0)
	repo.db.Model(&enrolments).Where("user_id", userId).Preload("Course").Find(&enrolments)
	return enrolments
}

func (repo *EnrolmentRepository) IsEnrolledInCourse(userId int32, courseId int64) *model.Enrolment {
	var enrolment model.Enrolment
	repo.db.Model(&enrolment).Where("user_id", userId).Where("course_id", courseId).Find(&enrolment)
	if enrolment.CourseID != 0 {
		return &enrolment
	}
	return nil
}
