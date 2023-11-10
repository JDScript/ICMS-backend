package message

import (
	"icms/internal/model"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (repo *MessageRepository) GetCourseMessages(courseId int64, userId int32) []model.CourseMessage {
	messages := make([]model.CourseMessage, 0)
	repo.db.Raw(`
	SELECT
		CM.id,
		CM.course_id,
		CM.title,
		CM.content,
		RM.read_at,
		CM.created_at,
		CM.updated_at
	FROM
		course_messages CM
	INNER JOIN enrolments E ON CM.course_id = e.course_id
	LEFT JOIN read_messages RM ON CM.id = RM.message_id AND E.user_id = RM.user_id
	WHERE
		E.user_id = ? AND E.course_id = ?
	ORDER BY CM.created_at DESC
	`, userId, courseId).Scan(&messages)
	return messages
}

func (repo *MessageRepository) GetUserMessages(userId int32) []model.CourseMessage {
	messages := make([]model.CourseMessage, 0)
	repo.db.Raw(`
	SELECT
		CM.id,
		CM.course_id,
		CM.title,
		CM.content,
		RM.read_at,
		CM.created_at,
		CM.updated_at
	FROM
		course_messages CM
	INNER JOIN enrolments E ON CM.course_id = e.course_id
	LEFT JOIN read_messages RM ON CM.id = RM.message_id AND E.user_id = RM.user_id
	WHERE
		E.user_id = ?
	ORDER BY CM.created_at DESC
	`, userId).Scan(&messages)
	return messages
}
