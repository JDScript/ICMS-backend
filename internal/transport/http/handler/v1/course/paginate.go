package course

import (
	"icms/internal/transport/http/request"
	"icms/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *CourseHandler) Paginate(c *gin.Context) {
	req := request.CoursePaginateRequest{}
	if ok := request.BindAndValidate(c, &req); !ok {
		return
	}

	paging := handler.courseRepo.PaginateCourses(c, &req)

	response.JSON(c, http.StatusOK, true, "Query success", paging)
}
