package me

import (
	"icms/pkg/auth"
	"icms/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *MeHandler) GetMessages(c *gin.Context) {
	user := auth.CurrentUser(c)

	messages := handler.messageRepo.GetUserMessages(user.ID)

	response.JSON(c, http.StatusOK, true, "Query success", messages)
}
