package model

import (
	"icms/internal/model/enum"
)

type User struct {
	ID    int32  `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"type:varchar(255)" json:"name"`
	Email string `gorm:"type:varchar(255);index:,unique;not null" json:"email"`

	CurrentLoginAt *int64 `gorm:"index" json:"current_login_at"`
	LastLoginAt    *int64 `gorm:"index" json:"last_login_at"`
	LastActivityAt *int64 `gorm:"index" json:"last_activity_at"`

	Identity    enum.UserIdentity  `gorm:"type:enum('student','teacher','admin');index;not null;default:student" json:"identity"`
	Descriptors []FacialDescriptor `gorm:"foreignKey:UserID"`

	CommonTimestampField
}
