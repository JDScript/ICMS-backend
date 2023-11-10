package bootstrap

import (
	appConfig "icms/internal/config"
	"icms/pkg/console"
	"icms/pkg/log"

	kratosConfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	kratosLog "github.com/go-kratos/kratos/v2/log"
	"github.com/thoas/go-funk"
)

var (
	config      kratosConfig.Config     // kratos config interface
	configModel = new(appConfig.Config) // app config model
)

func InitConfig(
	configPath string,
) *appConfig.Config {
	configResources := []kratosConfig.Source{file.NewSource(configPath)}

	kratosLog.SetLogger(log.NewEmptyLogger())

	config = kratosConfig.New(
		kratosConfig.WithSource(configResources...),
	)

	if err := config.Load(); err != nil {
		console.Exit("配置文件加载错误: " + err.Error())
	}

	if err := config.Scan(configModel); err != nil {
		console.Exit(err.Error())
	}

	// check that the environment is set correctly
	if !funk.ContainsString(appConfig.SupportedEnvs, configModel.App.Env.String()) {
		panic("unsupported env value: " + configModel.App.Env)
	}

	console.Success("当前配置环境为: " + configModel.App.Env.String())

	appConfig.GlobalConfig = configModel

	return configModel
}
