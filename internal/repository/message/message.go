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

	messages := make([]model.CourseMessage, 0)
	var total int64

	query := repo.db.
		Select(
			"CM.id",
			"CM.course_id",
			"CM.title",
			"CM.content",
			"CM.sender_id",
			"RM.read_at",
			"CM.created_at",
			"CM.updated_at",
			"U.id AS Sender__id",
			"U.name AS Sender__name",
			"U.email AS Sender__email",
		).
		Table("course_messages as CM").
		Joins("LEFT JOIN read_messages RM ON RM.user_id = ? AND CM.id = RM.message_id", userId).
		Joins("LEFT JOIN users U ON CM.sender_id = U.id").
		Where("CM.course_id", courseId)

	query.Count(&total)
	query.Order("CM.created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).Scan(&messages)

	return paginator.Paging{
		List:  messages,
		Total: total,
	}
}

func (repo *MessageRepository) GetUserMessages(c *gin.Context, userId int32, courseId *int64, unread bool) (paging paginator.Paging) {
	pageSize := cast.ToInt(c.DefaultQuery("page_size", "20"))
	page := cast.ToInt(c.DefaultQuery("page", "1"))

	messages := make([]model.CourseMessage, 0)
	var total int64

	query := repo.db.
		Select(
			"CM.id",
			"CM.course_id",
			"CM.title",
			"CM.content",
			"CM.sender_id",
			"RM.read_at",
			"CM.created_at",
			"CM.updated_at",
			"C.id AS Course__id",
			"C.code AS Course__code",
			"C.year AS Course__year",
			"C.section AS Course__section",
			"C.title AS Course__title",
			"C.summary AS Course__summary",
			"C.zoom_link AS Course__zoom_link",
			"U.id AS Sender__id",
			"U.name AS Sender__name",
			"U.email AS Sender__email",
		).
		Table("course_messages as CM").
		// Inner join 筛选用户要看到的 信息
		Joins("INNER JOIN enrolments E ON CM.course_id = E.course_id").
		Joins("LEFT JOIN read_messages RM ON CM.id = RM.message_id AND E.user_id = RM.user_id").
		Joins("LEFT JOIN courses C ON CM.course_id = C.id").
		Joins("LEFT JOIN users U ON CM.sender_id = U.id").
		Where("E.user_id", userId)

	if courseId != nil {
		query.Where("CM.course_id", courseId)
	}

	if unread {
		query.Where("RM.read_at IS NULL")
	}

	query.Count(&total)
	query.Order("CM.created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).Scan(&messages)

	return paginator.Paging{
		List:  messages,
		Total: total,
	}
}

func (repo *MessageRepository) ReadUserMessages(readMessages []model.ReadMessage) error {
	return repo.db.Create(&readMessages).Error
}

func (repo *MessageRepository) GetLatestCourseMessages(courseId int64, limit int64) []model.CourseMessage {
	latestMessages := []model.CourseMessage{}
	repo.db.
		Where("course_id", courseId).
		Order("updated_at DESC").
		Limit(3).
		Find(&latestMessages)

	return latestMessages
}
