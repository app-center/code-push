package zap_log

import (
	"github.com/funnyecho/code-push/pkg/log"
	"go.uber.org/zap"
)

func New(logger *zap.SugaredLogger) log.Logger {
	logger = logger.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar()
	return &loggerImpl{
		logger: logger,
	}
}

type loggerImpl struct {
	logger *zap.SugaredLogger
}

func (l *loggerImpl) Info(msg string, keyvals ...interface{}) {
	l.logger.Infow(msg, keyvals...)
}

func (l *loggerImpl) Debug(msg string, keyvals ...interface{}) {
	l.logger.Debugw(msg, keyvals...)
}
