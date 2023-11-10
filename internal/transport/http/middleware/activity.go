package middleware

import (
	"encoding/json"
	"icms/internal/model"
	"icms/internal/model/enum"
	"icms/internal/repository/activity"
	"icms/pkg/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type ActivityMiddleware func(actType enum.ActivityType) gin.HandlerFunc

func NewActivityMiddleware(activityRepo *activity.ActivityRepository) ActivityMiddleware {
	return func(actType enum.ActivityType) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			// Async
			go func() {
				// Remove sensitive fields
				ctx.Request.Header.Del("Authorization")
				ctx.Request.Header.Del("Cookie")

				headerBytes, err := json.Marshal(ctx.Request.Header)
				userId := auth.CurrentUID(ctx)

				if err == nil && userId != 0 {
					activity := model.Activity{
						Type:   actType,
						Path:   ctx.Request.RequestURI,
						Method: ctx.Request.Method,
						Header: datatypes.JSON(headerBytes),
						SrcIP:  ctx.ClientIP(),
						UserID: userId,
					}

					activityRepo.Create(&activity)
				}
			}()

			ctx.Next()
		}
	}
}
