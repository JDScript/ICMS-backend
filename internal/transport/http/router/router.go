package router

import (
	"time"

	"icms/internal/transport/http/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	New,
	NewV1Group,
)

func New(v1Group *v1Group) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 注册全局中间件
	router.Use(
		middleware.Recovery(),
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"*"},
			AllowHeaders:     []string{"*"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
	)

	// 加载 v1 路由
	v1Group.setup(router.Group("v1"))
	v1Group.useRoutes()

	return router
}
