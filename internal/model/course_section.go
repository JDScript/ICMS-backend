package model

type CourseSection struct {
	ID int64 `gorm:"primaryKey" json:"id"`

	CourseID int64  `gorm:"index;not null" json:"course_id"`
	Course   Course `gorm:"foreignKey:CourseID" json:"-"`

	Name    string `gorm:"varchar(255)" json:"name"`
	Summary string `gorm:"text" json:"summary"`

	Modules []CourseModule `gorm:"foreignKey:SectionID" json:"modules"`

	Order int64 `gorm:"index;not null;default:0" json:"order"`
	CommonTimestampField
}
