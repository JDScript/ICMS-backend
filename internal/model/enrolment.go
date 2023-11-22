package model

type Enrolment struct {
	CourseID int64  `gorm:"primaryKey" json:"course_id"`
	Course   Course `gorm:"foreignKey:CourseID" json:"course"`

	UserID int32 `gorm:"primaryKey" json:"user_id"`
	User   User  `gorm:"foreignKey:UserID" json:"-"`

	CommonTimestampField
}
