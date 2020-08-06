package log

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func New(logKit log.Logger) Logger {
	return &logger{
		Logger: logKit,
	}
}

type Logger interface {
	log.Logger
	Info(keyvals ...interface{}) error
	Debug(keyvals ...interface{}) error
}

type logger struct {
	log.Logger
}

func (l *logger) Info(keyvals ...interface{}) error {
	keyvals = append(keyvals, level.Key(), level.InfoValue())
	return l.Log(keyvals...)
}

func (l *logger) Debug(keyvals ...interface{}) error {
	keyvals = append(keyvals, level.Key(), level.DebugValue())
	return l.Log(keyvals...)
}
