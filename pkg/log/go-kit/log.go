package go_kig_log

import (
	log2 "github.com/funnyecho/code-push/pkg/log"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func New(logger log.Logger) log2.Logger {
	return &loggerImpl{
		logger: logger,
	}
}

type loggerImpl struct {
	logger log.Logger
}

func (l *loggerImpl) Info(msg string, keyvals ...interface{}) {
	keyvals = append([]interface{}{level.Key(), level.InfoValue(), "msg", msg}, keyvals...)
	l.logger.Log(keyvals...)
}

func (l *loggerImpl) Debug(msg string, keyvals ...interface{}) {
	keyvals = append([]interface{}{level.Key(), level.DebugValue(), "msg", msg}, keyvals...)
	l.logger.Log(keyvals...)
}
