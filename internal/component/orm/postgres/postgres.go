package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Driver               string
	Addr                 string
	Database             string
	Username             string
	Password             string
	Options              string
	MaxIdleConn          int
	MaxOpenConn          int
	ConnMaxIdleTime      time.Duration
	ConnMaxLifeTime      time.Duration
	Logger               logger.Interface
	Conn                 *sql.DB
	PreferSimpleProtocol bool
}

// New initialize *gorm.DB
func New(c Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  BuildDSN(c),
		PreferSimpleProtocol: c.PreferSimpleProtocol,
	}), &gorm.Config{
		SkipDefaultTransaction:                   false,
		DisableForeignKeyConstraintWhenMigrating: false,
		Logger:                                   c.Logger,
		QueryFields:                              true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if c.MaxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	}

	if c.MaxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	}

	if c.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(c.ConnMaxIdleTime)
	}

	if c.ConnMaxLifeTime > 0 {
		sqlDB.SetConnMaxLifetime(c.ConnMaxLifeTime)
	}

	return db, nil
}

// BuildDSN build dss to connect to the database
func BuildDSN(c Config) string {
	var host, port string
	s := strings.SplitN(c.Addr, ":", 2)
	if len(s) == 2 {
		host = s[0]
		port = s[1]
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", c.Username, c.Password, host, port, c.Database, c.Options)
	return dsn
}
