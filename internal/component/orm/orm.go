package orm

import (
	"context"
	"errors"
	"time"

	"icms/internal/component/orm/mysql"
	"icms/internal/component/orm/postgres"

	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var (
	ErrUnsupportedType         = errors.New("unsupported database type")
	ErrUnsupportedResolverType = errors.New("unsupported resolver type")
)

// Driver database driver type
type Driver string

func (d Driver) String() string {
	return string(d)
}

const (
	MySQL       Driver = "mysql"
	PostgresSQL Driver = "postgres"
)

// DSN database connection configuration
type DSN struct {
	Addr     string
	Database string
	Username string
	Password string
	Options  string
}

type Config struct {
	Driver Driver
	DSN
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxIdleTime time.Duration
	ConnMaxLifeTime time.Duration
	Resolvers       []Resolver
	QueryLogLevel   string
}

// New initialize orm
func New(config *Config, logger klog.Logger) (db *gorm.DB, cleanup func(), err error) {
	if config == nil {
		return nil, func() {}, nil
	}

	switch config.Driver {
	case MySQL:
		db, err = mysql.New(mysql.Config{
			Driver:                    config.Driver.String(),
			Addr:                      config.Addr,
			Database:                  config.Database,
			Username:                  config.Username,
			Password:                  config.Password,
			Options:                   config.Options,
			MaxIdleConn:               config.MaxIdleConn,
			MaxOpenConn:               config.MaxOpenConn,
			ConnMaxIdleTime:           config.ConnMaxIdleTime * time.Second,
			ConnMaxLifeTime:           config.ConnMaxLifeTime * time.Second,
			Conn:                      nil,
			SkipInitializeWithVersion: false,
			DefaultStringSize:         0,
			DisableDatetimePrecision:  false,
			DontSupportRenameIndex:    false,
			DontSupportRenameColumn:   false,
			Logger:                    &KratosLogWrapper{logger: logger, QueryLogLevel: klog.ParseLevel(config.QueryLogLevel)},
		})
		if err != nil {
			return
		}
	case PostgresSQL:
		db, err = postgres.New(postgres.Config{
			Driver:               config.Driver.String(),
			Addr:                 config.Addr,
			Database:             config.Database,
			Username:             config.Username,
			Password:             config.Password,
			Options:              config.Options,
			MaxIdleConn:          config.MaxIdleConn,
			MaxOpenConn:          config.MaxOpenConn,
			ConnMaxIdleTime:      config.ConnMaxIdleTime * time.Second,
			ConnMaxLifeTime:      config.ConnMaxLifeTime * time.Second,
			Conn:                 nil,
			PreferSimpleProtocol: false,
		})
		if err != nil {
			return
		}
	default:
		return nil, nil, ErrUnsupportedType
	}

	if len(config.Resolvers) > 0 {
		if err = registerResolver(db, config.Driver, config.Resolvers); err != nil {
			return nil, nil, err
		}
	}

	cleanup = func() {
		klog.NewHelper(logger).Info("closing the database resources")

		sqlDB, err := db.DB()
		if err != nil {
			klog.NewHelper(logger).Error(err)
		}

		if err := sqlDB.Close(); err != nil {
			klog.NewHelper(logger).Error(err)
		}
	}

	return db, cleanup, nil
}

func registerResolver(db *gorm.DB, driver Driver, resolvers []Resolver) error {
	if len(resolvers) > 0 {
		var (
			sources  = make([]gorm.Dialector, 0, len(resolvers))
			replicas = make([]gorm.Dialector, 0, len(resolvers))
		)

		for _, resolver := range resolvers {
			dial, err := BuildDialector(driver, resolver.DSN)
			if err != nil {
				return err
			}
			switch resolver.Type {
			case Source:
				sources = append(sources, dial)
			case Replica:
				replicas = append(replicas, dial)
			default:
				return ErrUnsupportedResolverType
			}
		}

		return db.Use(dbresolver.Register(dbresolver.Config{
			Sources:  sources,
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}))
	}

	return nil
}

type KratosLogWrapper struct {
	logger        klog.Logger
	QueryLogLevel klog.Level
}

func (l *KratosLogWrapper) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *KratosLogWrapper) Info(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Log(klog.LevelInfo, "msg", msg)
}

func (l *KratosLogWrapper) Warn(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Log(klog.LevelWarn, "msg", msg)

}
func (l *KratosLogWrapper) Error(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Log(klog.LevelError, "msg", msg)
}

func (l *KratosLogWrapper) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rowsAffected := fc()
	end := time.Now()
	elapsed := end.Sub(begin).Seconds()

	if err != nil {
		l.logger.Log(klog.LevelError, "msg", "[Database] "+err.Error(), "sql", sql, "rowsAffected", rowsAffected, "elapsed", elapsed)
	} else {
		l.logger.Log(l.QueryLogLevel, "msg", "[Database] Query", "sql", sql, "rowsAffected", rowsAffected, "elapsed", elapsed)
	}
}
