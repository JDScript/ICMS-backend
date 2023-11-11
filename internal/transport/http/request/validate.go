package request

import (
	"icms/pkg/response"

	"github.com/gin-gonic/gin"

	"github.com/bytedance/go-tagexpr/v2/binding"
)

type Error struct {
	Type      string `json:"type"`
	Msg       string `json:"msg"`
	FailField string `json:"fail_field"`
}

func (e *Error) Error() string {
	if e.Msg != "" {
		return e.Type + ": expr_path=" + e.FailField + ", cause=" + e.Msg
	}
	return e.Type + ": expr_path=" + e.FailField + ", cause=invalid"
}

// 这一部分进行了改造，使用了字节的 binder 和 validator
func BindAndValidate(c *gin.Context, obj interface{}) bool {
	binder := binding.New(nil)

	binder.SetErrorFactory(func(failField, msg string) error {
		return &Error{
			Type:      "binding",
			Msg:       msg,
			FailField: failField,
		}
	}, func(failField, msg string) error {
		return &Error{
			Type:      "validating",
			Msg:       msg,
			FailField: failField,
		}
	})
	// 解析并验证请求，文档见 https://github.com/bytedance/go-tagexpr
	if err := binder.BindAndValidate(obj, c.Request, c.Params); err != nil {
		response.ValidationError(c, err)
		return false
	}

	return true
}
