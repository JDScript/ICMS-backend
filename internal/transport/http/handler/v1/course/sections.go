package course

import (
	"icms/internal/transport/http/request"
	"icms/pkg/auth"
	"icms/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (hander *CourseHandler) GetSections(c *gin.Context) {
	req := request.CourseSectionsGetRequest{}
	if ok := request.BindAndValidate(c, &req); !ok {
		return
	}

	user := auth.CurrentUser(c)
	if enrol := hander.enrolmentRepo.IsEnrolledInCourse(user.ID, req.CourseID); enrol == nil {
		response.Abort403(c, "You haven't enrolled in this course")
		return
	}

	sections := hander.courseRepo.GetCourseContents(req.CourseID)

	response.JSON(c, http.StatusOK, true, "Query success", sections)
}
