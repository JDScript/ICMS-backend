package me

import (
	"icms/internal/model"
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

func (handler *MeHandler) ReadMessages(c *gin.Context) {
	req := request.MeReadMessagesRequest{}
	if ok := request.BindAndValidate(c, &req); !ok {
		return
	}

	user := auth.CurrentUser(c)

	readMessages := make([]model.ReadMessage, len(req.MessagesID))
	for idx := 0; idx < len(req.MessagesID); idx++ {
		readMessages[idx] = model.ReadMessage{
			MessageID: req.MessagesID[idx],
			UserID:    user.ID,
		}
	}

	err := handler.messageRepo.ReadUserMessages(readMessages)
	if err != nil {
		response.Abort500(c, err.Error())
		return
	}

	response.JSON(c, http.StatusNoContent, true, "Messages read", nil)
}
