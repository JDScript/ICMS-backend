package router

import (
	"icms/internal/model/enum"
	"icms/internal/transport/http/handler/v1/authentication"
	"icms/internal/transport/http/handler/v1/chat"
	"icms/internal/transport/http/handler/v1/course"
	"icms/internal/transport/http/handler/v1/face"
	"icms/internal/transport/http/handler/v1/me"
	"icms/internal/transport/http/handler/v1/user"
	"icms/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

type v1Group struct {
	authenticationHandler *authentication.AuthenticationHandler
	chatHandler           *chat.ChatHandler
	courseHandler         *course.CourseHandler
	faceHandler           *face.FaceHandler
	meHandler             *me.MeHandler
	userHandler           *user.UserHandler

	activityMiddleware middleware.ActivityMiddleware
	authMiddleware     middleware.JwtAuthMiddleware
	group              *gin.RouterGroup
}

func NewV1Group(
	authenticationHandler *authentication.AuthenticationHandler,
	chatHandler *chat.ChatHandler,
	courseHandler *course.CourseHandler,
	faceHandler *face.FaceHandler,
	meHandler *me.MeHandler,
	userHandler *user.UserHandler,
	activityMiddleware middleware.ActivityMiddleware,
	authMiddleware middleware.JwtAuthMiddleware,
) *v1Group {
	return &v1Group{
		authenticationHandler: authenticationHandler,
		chatHandler:           chatHandler,
		courseHandler:         courseHandler,
		faceHandler:           faceHandler,
		meHandler:             meHandler,
		userHandler:           userHandler,
		activityMiddleware:    activityMiddleware,
		authMiddleware:        authMiddleware,
	}
}

func (g *v1Group) setup(rg *gin.RouterGroup) {
	g.group = rg
}

// 具体路由
func (g *v1Group) useRoutes() {

	authenticationsGroup := g.group.Group("/authentications")
	{
		authenticationsGroup.POST("", g.authenticationHandler.Login, g.activityMiddleware(enum.Activity_Login))
	}

	meGroup := g.group.Group("/me", g.authMiddleware())
	{
		meGroup.GET("", g.meHandler.GetMe, g.activityMiddleware(enum.Activity_Get_My_Profile))
		meGroup.GET("/activities", g.meHandler.GetActivities, g.activityMiddleware(enum.Activity_Get_My_Activities))
		meGroup.DELETE("/activities", g.meHandler.ClearActivities, g.activityMiddleware(enum.Activity_Clear_My_Activities))
		meGroup.GET("/enrolments", g.meHandler.GetEnrolments, g.activityMiddleware(enum.Activity_Get_My_Enrolments))
		meGroup.POST("/enrolments", g.meHandler.CreateEnrolment, g.activityMiddleware(enum.Activity_Greate_My_Enrolment))
		meGroup.GET("/messages", g.meHandler.GetMessages, g.activityMiddleware(enum.Activity_Get_My_Messages))
		meGroup.DELETE("/messages", g.meHandler.ReadMessages, g.activityMiddleware(enum.Activity_Read_My_Messages))
	}

	userGroup := g.group.Group("/users")
	{
		userGroup.POST("", g.userHandler.Create, g.activityMiddleware(enum.Activity_Register))
	}

	faceGroup := g.group.Group("/faces")
	{
		faceGroup.POST("", g.faceHandler.Detect)
	}

	courseGroup := g.group.Group("/courses", g.authMiddleware())
	{
		courseGroup.GET("", g.courseHandler.Paginate, g.activityMiddleware(enum.Activity_Search_All_Courses))
		courseGroup.GET("/:courseId", g.courseHandler.Get, g.activityMiddleware(enum.Activity_Get_Course_Detail))
		courseGroup.GET("/:courseId/sections", g.courseHandler.GetSections, g.activityMiddleware(enum.Activity_Get_Course_Sections))
		courseGroup.GET("/:courseId/messages", g.courseHandler.GetMessages, g.activityMiddleware(enum.Activity_Get_Course_Messages))
	}

	chatGroup := g.group.Group("/chat")
	{
		chatGroup.Any("/*path", g.authMiddleware(), g.chatHandler.ChatCompletions)
	}

	g.group.GET("/ms/refresh_token", g.chatHandler.RefreshToken)
}
