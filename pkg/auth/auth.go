package auth

import (
	"icms/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func CurrentUID(c *gin.Context) int32 {
	userId := cast.ToInt32(c.GetString("current_user_id"))

	return userId
}

func CurrentUser(c *gin.Context) *model.User {
	user, ok := c.MustGet("current_user").(*model.User)
	if !ok {
		return &model.User{}
	}

	return user
}
