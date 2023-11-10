package command

import (
	"icms/internal/component"
	"icms/internal/config"

	"github.com/google/wire"
	"github.com/spf13/cobra"
)

var ProviderSet = wire.NewSet(
	config.ProviderSet,
	component.ProviderSet,
	wire.NewSet(NewMigrateCommand),
)

type AllCommands struct {
	migrate MigrateCommand
}

func Setup(rootCommand *cobra.Command, newCommand func() (*AllCommands, func(), error)) {
	rootCommand.AddCommand(
		&cobra.Command{
			Use:   "migrate",
			Short: "迁移数据库",
			Run: func(cmd *cobra.Command, args []string) {
				command, cleanup, err := newCommand()
				if err != nil {
					panic(err)
				}
				defer cleanup()
				command.migrate.Migrate(cmd, args)
			},
		},
	)
}

func New(migrate MigrateCommand) *AllCommands {
	return &AllCommands{
		migrate: migrate,
	}
}
