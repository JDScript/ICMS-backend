package middleware

import (
	"net"
	"os"
	"strings"

	"icms/pkg/response"

	"github.com/gin-gonic/gin"
)

// Recovery 使用 zap.Error() 来记录 Panic 和 call stack
func Recovery() gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				// 链接中断，客户端中断连接为正常行为，不需要记录堆栈信息
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						errStr := strings.ToLower(se.Error())
						if strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				// 链接中断的情况
				if brokenPipe {
					// 链接已断开，无法写状态码
					return
				}

				// 返回 500 状态码
				response.Abort500(c)
			}
		}()
		c.Next()
	}
}
