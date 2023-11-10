package handler

import (
	"icms/internal/transport/http/handler/v1/authentication"
	"icms/internal/transport/http/handler/v1/course"
	"icms/internal/transport/http/handler/v1/face"
	"icms/internal/transport/http/handler/v1/me"
	"icms/internal/transport/http/handler/v1/user"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	authentication.NewHandler,
	course.NewHandler,
	face.NewHandler,
	me.NewHandler,
	user.NewHandler,
)
