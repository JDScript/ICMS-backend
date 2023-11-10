package http

import (
	"fmt"

	"icms/internal/config"
	"icms/internal/transport/http/handler"
	"icms/internal/transport/http/middleware"
	"icms/internal/transport/http/router"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
)

// @securityDefinitions.apiKey	Bearer
// @in							header
// @name						Authorization
// @scheme						bearer
var ProviderSet = wire.NewSet(
	handler.ProviderSet,
	middleware.ProviderSet,
	router.ProviderSet,
	NewHTTPServer,
)

func NewHTTPServer(config *config.Config, g *gin.Engine) *http.Server {
	srv := http.NewServer(http.Address(fmt.Sprintf("%s:%d", config.HTTP.Host, config.HTTP.Port)))
	srv.HandlePrefix("/", g)

	return srv
}
