package model

type CourseMessage struct {
	ID       int64   `gorm:"primaryKey" json:"id"`
	CourseID int64   `gorm:"index" json:"course_id"`
	Course   *Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`

	Title   string `gorm:"varchar(255);not null" json:"title"`
	Content string `gorm:"text;not null" json:"content"`

	SenderID int32 `gorm:"index;not null" json:"sender_id"`
	Sender   User  `gorm:"foreignKey:SenderID" json:"sender,omitempty"`

	ReadAt *int64 `gorm:"-:migration" json:"read_at"`

	CommonTimestampField
}
