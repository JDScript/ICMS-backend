package request

import (
	"icms/pkg/response"

	"github.com/gin-gonic/gin"

	"github.com/bytedance/go-tagexpr/v2/binding"
)

// 这一部分进行了改造，使用了字节的 binder 和 validator
func BindAndValidate(c *gin.Context, obj interface{}) bool {

	binder := binding.New(nil)
	// 解析并验证请求，文档见 https://github.com/bytedance/go-tagexpr
	if err := binder.BindAndValidate(obj, c.Request, c.Params); err != nil {
		response.ValidationError(c, err)
		return false
	}

	return true
}
