package course

import (
	"icms/internal/model"
	"icms/internal/transport/http/request"
	"icms/pkg/paginator"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *CourseRepository {
	return &CourseRepository{
		db: db,
	}
}

func (repo *CourseRepository) PaginateCourses(
	ctx *gin.Context,
	req *request.CoursePaginateRequest,
) (paging paginator.Paging) {
	courses := []model.Course{}
	query := repo.db.Model(&courses)

	if req.Search != nil {
		q := strings.ToLower("%" + *req.Search + "%")
		query = query.Where("LOWER(code) LIKE ?", q).Or("LOWER(title) LIKE ?", q)
	}

	paging = paginator.Paginate(
		ctx,
		query,
		&courses,
	)

	return
}

func (repo *CourseRepository) GetCourseContents(courseId int64) []model.CourseSection {
	sections := make([]model.CourseSection, 0)
	repo.db.Model(&sections).Where("course_id", courseId).Preload("Modules").Find(&sections)
	return sections
}

func (repo *CourseRepository) GetCourse(courseId int64) (course *model.Course) {
	repo.db.Where("id", courseId).Find(&course)
	return
}
