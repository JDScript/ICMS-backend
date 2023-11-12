package course

import (
	"icms/internal/transport/http/request"
	"icms/pkg/auth"
	"icms/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (hander *CourseHandler) GetMessages(c *gin.Context) {
	req := request.CourseMessagesGetRequest{}
	if ok := request.BindAndValidate(c, &req); !ok {
		return
	}

	user := auth.CurrentUser(c)
	if enrol := hander.enrolmentRepo.IsEnrolledInCourse(user.ID, req.CourseID); enrol == nil {
		response.Abort403(c, "You haven't enrolled in this course")
		return
	}

	paging := hander.messageRepo.GetCourseMessages(c, req.CourseID, user.ID)

	response.JSON(c, http.StatusOK, true, "Query success", paging)
}
