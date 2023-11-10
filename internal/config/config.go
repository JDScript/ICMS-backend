package config

import (
	"icms/internal/component/facial"
	"icms/internal/component/jwt"
	"icms/internal/component/orm"

	"github.com/google/wire"
)

var GlobalConfig *Config

var ProviderSet = wire.NewSet(
	wire.FieldsOf(
		new(*Config),
		"App",
		"HTTP",
		"DB",
		"Jwt",
		"Facial",
	),
)

// SupportedEnvs 支持的环境
var SupportedEnvs = []string{Local.String(), Staging.String(), Prod.String()}

// Env 当前运行环境，当前设置 local, staging 和 prod
type Env string

func (e Env) String() string {
	return string(e)
}

const (
	Local   Env = "local"
	Staging Env = "staging"
	Prod    Env = "prod"
)

type Config struct {
	App    *App           `json:"app"`
	HTTP   *HTTP          `json:"http"`
	DB     *orm.Config    `json:"db"`
	Jwt    *jwt.Config    `json:"jwt"`
	Facial *facial.Config `json:"facial"`
}

type App struct {
	Name string `json:"name"`
	Env  Env    `json:"env"`
}

type HTTP struct {
	Host string `json:"host"`
	Port int64  `json:"port"`
}
