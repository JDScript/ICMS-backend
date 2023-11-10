package transport

import (
	"icms/internal/transport/http"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	http.ProviderSet,
)
