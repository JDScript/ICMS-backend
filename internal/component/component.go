package component

import (
	"icms/internal/component/chatgpt"
	"icms/internal/component/facial"
	"icms/internal/component/jwt"
	"icms/internal/component/orm"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(chatgpt.New),
	wire.NewSet(jwt.New),
	wire.NewSet(orm.New),
	wire.NewSet(facial.New),
)
