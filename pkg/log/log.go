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
	Info(msg string, keyvals ...interface{}) error
	Debug(msg string, keyvals ...interface{}) error
}

type logger struct {
	log.Logger
}

func (l *logger) Info(msg string, keyvals ...interface{}) error {
	keyvals = append([]interface{}{level.Key(), level.InfoValue(), "msg", msg}, keyvals...)
	return l.Log(keyvals...)
}

func (l *logger) Debug(msg string, keyvals ...interface{}) error {
	keyvals = append([]interface{}{level.Key(), level.DebugValue(), "msg", msg}, keyvals...)
	return l.Log(keyvals...)
}
