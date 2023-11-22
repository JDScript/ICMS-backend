package activity

import (
	"icms/internal/model"
	"icms/internal/transport/http/request"
	"icms/pkg/paginator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ActivityRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{
		db: db,
	}
}

func (repo *ActivityRepository) Create(act *model.Activity) error {
	return repo.db.Create(act).Error
}

func (repo *ActivityRepository) PaginateByUser(
	c *gin.Context,
	userId int32,
	req *request.ActivityPaginateRequest,
) (paging paginator.Paging) {
	activities := []model.Activity{}
	query := repo.db.Model(&activities).Where("user_id", userId)

	if req.Method != nil && len(*req.Method) > 0 {
		query = query.Where("method IN ?", *req.Method)
	}

	if req.Type != nil && len(*req.Type) > 0 {
		query = query.Where("type IN ?", *req.Type)
	}

	paging = paginator.Paginate(
		c,
		query,
		&activities,
	)
	return
}

func (repo *ActivityRepository) Clear(userId int32) error {
	return repo.db.Where("user_id", userId).Delete(&model.Activity{}).Error
}
