package internal

type InterceptorLogger interface {
	Info(msg string, keyvals ...interface{}) error
}
