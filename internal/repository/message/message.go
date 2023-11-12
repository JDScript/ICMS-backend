package message

import (
	"icms/internal/model"
	"icms/pkg/paginator"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
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

func (repo *MessageRepository) GetCourseMessages(c *gin.Context, courseId int64, userId int32) (paging paginator.Paging) {
	pageSize := cast.ToInt(c.DefaultQuery("page_size", "20"))
	page := cast.ToInt(c.DefaultQuery("page", "1"))

	messages := make([]model.CourseMessage, pageSize)
	var total int64

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
	LIMIT ? OFFSET ?
	`, userId, courseId, pageSize, (page-1)*pageSize)

	repo.db.Raw(`
	SELECT
		COUNT(CM.id)
	FROM
		course_messages CM
	INNER JOIN enrolments E ON CM.course_id = E.course_id
	LEFT JOIN read_messages RM ON CM.id = RM.message_id AND E.user_id = RM.user_id
	WHERE
		E.user_id = ? AND E.course_id = ?
	GROUP BY CM.id
	`, userId, courseId).Scan(&total)

	return paginator.Paging{
		List:  messages,
		Total: total,
	}
}

func (repo *MessageRepository) GetUserMessages(c *gin.Context, userId int32) (paging paginator.Paging) {
	pageSize := cast.ToInt(c.DefaultQuery("page_size", "20"))
	page := cast.ToInt(c.DefaultQuery("page", "1"))

	messages := make([]model.CourseMessage, 0)
	var total int64

	repo.db.Raw(`
	SELECT
		CM.id,
		CM.course_id,
		CM.title,
		CM.content,
		RM.read_at,
		CM.created_at,
		CM.updated_at,
		C.id AS Course__id,
		C.code AS Course__code,
		C.year AS Course__year,
		C.section AS Course__section,
		C.title AS Course__title,
		C.instructor AS Course__instructor,
		C.summary AS Course__summary,
		C.zoom_link AS Course__zoom_link
	FROM
		course_messages CM
	INNER JOIN enrolments E ON CM.course_id = E.course_id
	LEFT JOIN read_messages RM ON CM.id = RM.message_id AND E.user_id = RM.user_id
	LEFT JOIN courses C ON CM.course_id = C.id
	WHERE
		E.user_id = ?
	ORDER BY CM.created_at DESC
	LIMIT ? OFFSET ?
	`, userId, pageSize, (page-1)*pageSize).Scan(&messages)

	repo.db.Raw(`
	SELECT
		COUNT(CM.id)
	FROM
		course_messages CM
	INNER JOIN enrolments E ON CM.course_id = E.course_id
	WHERE
		E.user_id = ?
	GROUP BY CM.id
	`, userId).Scan(&total)

	return paginator.Paging{
		List:  messages,
		Total: total,
	}
}

func (repo *MessageRepository) GetUserUnreadMessages(c *gin.Context, userId int32) (paging paginator.Paging) {
	pageSize := cast.ToInt(c.DefaultQuery("page_size", "20"))
	page := cast.ToInt(c.DefaultQuery("page", "1"))

	messages := make([]model.CourseMessage, 0)
	var total int64

	repo.db.Raw(`
	SELECT
		CM.id,
		CM.course_id,
		CM.title,
		CM.content,
		RM.read_at,
		CM.created_at,
		CM.updated_at,
		C.id AS Course__id,
		C.code AS Course__code,
		C.year AS Course__year,
		C.section AS Course__section,
		C.title AS Course__title,
		C.instructor AS Course__instructor,
		C.summary AS Course__summary,
		C.zoom_link AS Course__zoom_link
	FROM
		course_messages CM
	INNER JOIN enrolments E ON CM.course_id = E.course_id
	LEFT JOIN read_messages RM ON CM.id = RM.message_id AND E.user_id = RM.user_id
	LEFT JOIN courses C ON CM.course_id = C.id
	WHERE
		E.user_id = ? AND RM.read_at IS NULL
	ORDER BY CM.created_at DESC
	LIMIT ? OFFSET ?
	`, userId, pageSize, (page-1)*pageSize).Scan(&messages)

	repo.db.Raw(`
	SELECT
		COUNT(CM.id)
	FROM
		course_messages CM
	INNER JOIN enrolments E ON CM.course_id = E.course_id
	LEFT JOIN read_messages RM ON CM.id = RM.message_id AND E.user_id = RM.user_id
	WHERE
		E.user_id = ? AND RM.read_at IS NULL
	GROUP BY CM.id
	`, userId).Scan(&total)

	return paginator.Paging{
		List:  messages,
		Total: total,
	}
}
