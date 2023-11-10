package model

import "icms/internal/model/enum"

type Enrolment struct {
	CourseID int64  `gorm:"primaryKey;index" json:"course_id"`
	Course   Course `gorm:"foreignKey:CourseID" json:"course"`

	UserID int32 `gorm:"primaryKey;index" json:"user_id"`
	User   User  `gorm:"foreignKey:UserID" json:"-"`

	Identity enum.EnrolmentIdentity `gorm:"type:enum('student','assistant','teacher');index;not null;default:student" json:"identity"`

	CommonTimestampField
}
