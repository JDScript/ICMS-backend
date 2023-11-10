package log

import (
	"github.com/go-kratos/kratos/v2/log"
)

var _ log.Logger = (*EmptyLogger)(nil)

type EmptyLogger struct {
}

func NewEmptyLogger() *EmptyLogger {
	return &EmptyLogger{}
}

func (l *EmptyLogger) Log(level log.Level, keyvals ...interface{}) error {
	return nil
}

func (l *EmptyLogger) Sync() error {
	return nil
}

func (l *EmptyLogger) Close() error {
	return nil
}
