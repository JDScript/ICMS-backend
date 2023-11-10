package main

import (
	"context"
	"encoding/gob"
	"icms/internal/bootstrap"
	"icms/internal/command"
	"icms/internal/config"
	"icms/pkg/console"
	"icms/pkg/path"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Kagami/go-face"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	rootPath = path.RootPath()
)

func newApp(cfg *config.Config, logger log.Logger, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.Name(cfg.App.Name),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
		),
	)
}

func main() {
	pflag.Parse()
	// 初始化配置文件
	config := bootstrap.InitConfig(
		filepath.Join(rootPath, "config"),
	)

	gob.Register(face.Descriptor{})

	logger, zapLogger := bootstrap.SetupLogger(config)

	rootCommand := &cobra.Command{
		Use: "app",
		Run: func(cmd *cobra.Command, args []string) {
			app, cleanup, err := initApp(config, logger, zapLogger)
			if err != nil {
				console.Exit(err.Error())
			}
			defer cleanup()

			if err = app.Run(); err != nil {
				console.Exit(err.Error())
			}

			// monitor signal
			signalCtx, signalStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
			defer signalStop()

			// waiting exit signal ...
			<-signalCtx.Done()
			signalStop()

			// stop app
			app.Stop()
		},
	}

	command.Setup(rootCommand, func() (*command.AllCommands, func(), error) {
		return initCommand(config, logger, zapLogger)
	})

	if err := rootCommand.Execute(); err != nil {
		panic(err)
	}
}
