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
