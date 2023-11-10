package command

import (
	"os"

	"icms/internal/model"
	"icms/pkg/console"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

type MigrateCommand struct {
	db     *gorm.DB
	logger *log.Helper
}

func (m *MigrateCommand) Migrate(cmd *cobra.Command, args []string) {
	if m.db != nil {
		err := m.db.AutoMigrate(
			model.Activity{},
			model.CourseMessage{},
			model.CourseModule{},
			model.CourseSection{},
			model.Course{},
			model.User{},
			model.FacialDescriptor{},
			model.Enrolment{},
			model.ReadMessage{},
		)
		if err != nil {
			m.logger.Error(err)
			console.Error(err.Error())
		}
		os.Exit(0)
	}
}

func NewMigrateCommand(db *gorm.DB, logger log.Logger) MigrateCommand {
	return MigrateCommand{
		db:     db,
		logger: log.NewHelper(logger),
	}
}
