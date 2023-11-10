//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"go.uber.org/zap"
	"icms/internal/command"
	"icms/internal/component"
	"icms/internal/config"
	"icms/internal/repository"
	"icms/internal/transport"
)

func initApp(*config.Config, log.Logger, *zap.Logger) (*kratos.App, func(), error) {
	panic(
		wire.Build(
			component.ProviderSet,
			config.ProviderSet,
			repository.ProviderSet,
			transport.ProviderSet,
			newApp,
		),
	)
}

func initCommand(*config.Config, log.Logger, *zap.Logger) (*command.AllCommands, func(), error) {
	panic(wire.Build(
		command.ProviderSet,
		command.New,
	))
}
