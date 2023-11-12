package me

import (
	"icms/internal/transport/http/request"
	"icms/pkg/auth"
	"icms/pkg/paginator"
	"icms/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *MeHandler) GetMessages(c *gin.Context) {
	req := request.MeGetMessagesRequest{}
	if ok := request.BindAndValidate(c, &req); !ok {
		return
	}

	user := auth.CurrentUser(c)
	var paging paginator.Paging

	if req.Unread {
		paging = handler.messageRepo.GetUserUnreadMessages(c, user.ID)
	} else {
		paging = handler.messageRepo.GetUserMessages(c, user.ID)
	}

	response.JSON(c, http.StatusOK, true, "Query success", paging)
}
