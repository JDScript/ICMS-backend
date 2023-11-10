package middleware

import (
	"icms/internal/component/jwt"
	"icms/internal/repository/user"
	"icms/pkg/response"

	"github.com/gin-gonic/gin"
)

type JwtAuthMiddleware func() gin.HandlerFunc

func NewJWTMiddleware(jwt *jwt.Jwt, userRepo *user.UserRepository) JwtAuthMiddleware {
	return func() gin.HandlerFunc {
		return func(ctx *gin.Context) {
			claims, err := jwt.ParserToken(ctx)

			// JWT 解析失败，有错误发生
			if err != nil {
				response.Unauthorized(ctx, "Unauthenticated")
				return
			}

			// 获取用户
			user := userRepo.GetByID(claims.UserId)
			if user.ID == 0 {
				response.Unauthorized(ctx, "User not exists")
				return
			}

			// 将用户信息存入 gin.context 里，后续 auth 包将从这里拿到当前用户数据
			ctx.Set("current_user_id", claims.UserId)
			ctx.Set("current_user", user)
			ctx.Next()
		}
	}
}
