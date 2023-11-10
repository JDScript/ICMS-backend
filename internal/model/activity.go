package model

import (
	"icms/internal/model/enum"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Activity struct {
	ID   int64             `gorm:"primaryKey" json:"id"`
	Type enum.ActivityType `gorm:"index;not null;default:unknown" json:"type"`

	Path   string         `gorm:"varchar(255);not null" json:"path"`
	Method string         `gorm:"varchar(10);not null" json:"method"`
	Header datatypes.JSON `json:"header"`
	SrcIP  string         `gorm:"varchar(64);not null" json:"src_ip"`

	UserID int32 `gorm:"index;not null" json:"-"`
	User   User  `gorm:"foreignKey:UserID" json:"-"`

	CommonTimestampField
}

func (act *Activity) AfterSave(tx *gorm.DB) (err error) {
	// Update user's field of last login and last activity
	fields := map[string]interface{}{
		"last_activity_at": act.CommonTimestampField.CreatedAt,
	}

	if act.Type == enum.Activity_Login {
		fields["last_login_at"] = tx.Raw("SELECT a.current_login_at FROM (SELECT current_login_at FROM users WHERE id = ?) a", act.UserID)
		fields["current_login_at"] = act.CommonTimestampField.CreatedAt
	}

	err = tx.Table("users").Where("id", act.UserID).UpdateColumns(fields).Error
	if err != nil {
		return
	}

	return
}
