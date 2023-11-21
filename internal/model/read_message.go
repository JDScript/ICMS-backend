package model

type ReadMessage struct {
	MessageID int64         `gorm:"primaryKey"`
	Message   CourseMessage `gorm:"foreignKey:MessageID" json:"-"`
	UserID    int32         `gorm:"primaryKey"`
	User      User          `gorm:"foreignKey:UserID" json:"-"`

	ReadAt int64 `gorm:"autoCreateTime:milli" json:"created_at"`
}
