// Package response 响应处理工具
package response

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Response struct {
	Success bool        `json:"success"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSON(c *gin.Context, status int, success bool, message string, data interface{}) {
	statusText := http.StatusText(status)
	var msg string
	if len(message) > 0 {
		msg = message
	} else {
		msg = strconv.Itoa(status) + " " + statusText
	}

	c.JSON(status, Response{
		Success: success,
		Status:  status,
		Message: msg,
		Data:    data,
	})
}

func AbortJSON(c *gin.Context, status int, success bool, message string, data interface{}) {
	statusText := http.StatusText(status)
	var msg string
	if len(message) > 0 {
		msg = message
	} else {
		msg = strconv.Itoa(status) + " " + statusText
	}
	c.AbortWithStatusJSON(status, Response{
		Success: success,
		Status:  status,
		Message: msg,
		Data:    data,
	})
}

// Abort404 响应 404，未传参 msg 时使用默认消息
func Abort404(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusNotFound, Response{
		Success: false,
		Status:  http.StatusNotFound,
		Message: defaultMessage("Not found", msg...),
	})
}

// Abort403 响应 403，未传参 msg 时使用默认消息
func Abort403(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusForbidden, Response{
		Success: false,
		Status:  http.StatusForbidden,
		Message: defaultMessage("没有权限进行此操作", msg...),
	})
}

// Abort500 响应 500，未传参 msg 时使用默认消息
func Abort500(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
		Success: false,
		Status:  http.StatusInternalServerError,
		Message: defaultMessage("Server internal error", msg...),
	})
}

// BadRequest 响应 400，传参 err 对象，未传参 msg 时使用默认消息
// 在解析用户请求，请求的格式或者方法不符合预期时调用
func BadRequest(c *gin.Context, err error, msg ...string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, Response{
		Success: false,
		Status:  http.StatusBadRequest,
		Message: defaultMessage("Bad request", msg...),
		Data:    err.Error(),
	})
}

// Error 响应 404 或 422，未传参 msg 时使用默认消息
// 处理请求时出现错误 err，会附带返回 error 信息，如登录错误、找不到 ID 对应的 Model
func Error(c *gin.Context, err error, msg ...string) {
	// error 类型为『数据库未找到内容』
	if err == gorm.ErrRecordNotFound {
		Abort404(c)
		return
	}

	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, Response{
		Success: false,
		Status:  http.StatusUnprocessableEntity,
		Message: defaultMessage("Receieved unprocessable request", msg...),
		Data:    err.Error(),
	})
}

// ValidationError 处理表单验证不通过的错误，返回的 JSON 示例：
func ValidationError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, Response{
		Success: false,
		Status:  http.StatusUnprocessableEntity,
		Message: "Request validation error",
		Data:    err,
	})
}

// Unauthorized 响应 401，未传参 msg 时使用默认消息
// 登录失败、jwt 解析失败时调用
func Unauthorized(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
		Success: false,
		Status:  http.StatusUnauthorized,
		Message: defaultMessage("Unauthenticated", msg...),
	})
}

// defaultMessage 内用的辅助函数，用以支持默认参数默认值
// Go 不支持参数默认值，只能使用多变参数来实现类似效果
func defaultMessage(defaultMsg string, msg ...string) (message string) {
	if len(msg) > 0 {
		message = msg[0]
	} else {
		message = defaultMsg
	}
	return
}
