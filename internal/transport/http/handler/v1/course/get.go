package course

import (
	"icms/internal/transport/http/request"
	"icms/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *CourseHandler) Get(c *gin.Context) {
	req := request.CourseGetRequest{}
	if ok := request.BindAndValidate(c, &req); !ok {
		return
	}

	course := handler.courseRepo.GetCourse(req.CourseID)

	if course.ID == 0 {
		response.Abort404(c, "Course not exists")
		return
	}

	response.JSON(c, http.StatusOK, true, "Query success", course)
}
