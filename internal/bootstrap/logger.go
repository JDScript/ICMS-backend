package bootstrap

import (
	"io"
	"os"
	"path/filepath"

	appConfig "icms/internal/config"
	"icms/pkg/console"
	"icms/pkg/log"
	"icms/pkg/path"

	kratosZap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	kratosLog "github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

func SetupLogger(config *appConfig.Config) (kratosLog.Logger, *zap.Logger) {
	var writer io.Writer
	logLevel := "debug"

	logPath := filepath.Join(path.RootPath(), "logs/log_%Y%m%d.log")

	if config.App.Env == appConfig.Local {
		writer = os.Stdout
	} else {
		loggerWriter, err := rotatelogs.New(
			logPath,
			rotatelogs.WithClock(rotatelogs.Local),
		)
		if err != nil {
			console.Exit(err.Error())
		}
		writer = loggerWriter
		logLevel = "info"
	}

	zapLogger := log.New(
		log.WithLevel(log.Level(logLevel)),
		log.WithFormat(log.Format("json")),
		log.WithWriter(writer),
	)

	logger := kratosLog.With(
		kratosZap.NewLogger(zapLogger),
	)

	kratosLog.SetLogger(logger)

	return logger, zapLogger
}
