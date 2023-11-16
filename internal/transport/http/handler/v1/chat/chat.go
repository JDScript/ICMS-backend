package chat

import (
	chat "icms/internal/component/chatgpt"
	"icms/pkg/response"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chat *chat.ChatGPT
}

func NewHandler(chat *chat.ChatGPT) *ChatHandler {
	return &ChatHandler{
		chat: chat,
	}
}

func (h *ChatHandler) ChatCompletions(c *gin.Context) {
	remote := &h.chat.Endpoint

	c.Writer.Header().Del("Access-Control-Allow-Origin")

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = remote.Path + c.Param("path")
		req.Header = http.Header{
			"Authorization": {"Bearer " + h.chat.Resp.AccessToken},
		}
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

func (h *ChatHandler) RefreshToken(c *gin.Context) {
	err := h.chat.Refresh()
	response.JSON(c, http.StatusOK, true, "", err)
}
