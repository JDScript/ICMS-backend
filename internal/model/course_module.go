package model

import (
	"icms/internal/model/enum"

	"gorm.io/datatypes"
)

type CourseModule struct {
	ID int64 `gorm:"primaryKey" json:"id"`

	CourseID  int64         `gorm:"index;not null" json:"course_id"`
	Course    Course        `gorm:"foreignKey:CourseID" json:"-"`
	SectionID int64         `gorm:"index;not null" json:"section_id"`
	Section   CourseSection `gorm:"foreignKey:SectionID" json:"-"`

	Name       string                `gorm:"varchar(255);not null" json:"name"`
	ModuleType enum.CourseModuleType `gorm:"varchar(255);not null" json:"module_type"`
	Indent     int8                  `gorm:"default:0" json:"indent"`
	Link       string                `gorm:"varchar(255)" json:"link"`
	Extra      datatypes.JSON        `json:"extra"`

	Order int64 `gorm:"index;not null;default:0" json:"order"`
	CommonTimestampField
}
