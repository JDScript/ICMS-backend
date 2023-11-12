package me

import (
	"icms/internal/transport/http/request"
	"icms/pkg/auth"
	"icms/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *MeHandler) GetActivities(c *gin.Context) {
	req := request.ActivityPaginateRequest{}
	if ok := request.BindAndValidate(c, &req); !ok {
		return
	}

	userId := auth.CurrentUID(c)
	paging := handler.activityRepo.PaginateByUser(c, userId, &req)

	response.JSON(c, http.StatusOK, true, "Query success", paging)
}

func (handler *MeHandler) ClearActivities(c *gin.Context) {
	userId := auth.CurrentUID(c)
	err := handler.activityRepo.Clear(userId)
	if err != nil {
		response.Abort500(c, err.Error())
		return
	}

	response.JSON(c, http.StatusNoContent, true, "Cleared all activities", nil)
}
