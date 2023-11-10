package model

import (
	"database/sql/driver"
	"encoding/json"
)

type Course struct {
	ID         int64   `gorm:"primaryKey" json:"id"`
	Code       string  `gorm:"varchar(16);index:idx_course_uniq,unique;index;not null" json:"code"`
	Year       int     `gorm:"index:idx_course_uniq,unique;not null" json:"year"`
	Section    string  `gorm:"varchar(4);index:idx_course_uniq,unique;not null" json:"section"`
	Title      string  `gorm:"varchar(255);index;not null" json:"title"`
	Instructor string  `gorm:"varchar(255);" json:"instructor"`
	Summary    string  `gorm:"text" json:"summary"`
	ZoomLink   *string `gorm:"varchar(255)" json:"zoom_link"`

	Timeslots CourseTimeslots `gorm:"type:json" json:"slots"`

	CommonTimestampField
}

type CourseTimeslot struct {
	Day       uint8  `json:"day"`
	Venue     string `json:"venue"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Remark    string `json:"remark"`
}

type CourseTimeslots []CourseTimeslot

func (p CourseTimeslots) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *CourseTimeslots) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), p)
}
