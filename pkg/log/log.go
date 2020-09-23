package log

import "fmt"

type Logger interface {
	Info(msg string, keyvals ...interface{})
	Debug(msg string, keyvals ...interface{})
}

var NoopLogger = &noopLogger{}

type noopLogger struct{}

func (l *noopLogger) Info(msg string, keyvals ...interface{}) {}

func (l *noopLogger) Debug(msg string, keyvals ...interface{}) {}

var PanicLogger = &panicLogger{}

type panicLogger struct{}

func (l *panicLogger) Info(msg string, keyvals ...interface{}) {
	panic(fmt.Errorf("log to panic: %v", append(keyvals, "level", "info", "msg", msg)))
}

func (l *panicLogger) Debug(msg string, keyvals ...interface{}) {
	panic(fmt.Errorf("log to panic: %v", append(keyvals, "level", "debug", "msg", msg)))
}
