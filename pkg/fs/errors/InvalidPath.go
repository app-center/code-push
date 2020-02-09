package errors

import "github.com/funnyecho/code-push/pkg/errors"

type InvalidPathError error

type InvalidPathConfig struct {
	Err  error
	Path interface{}
}

func NewInvalidPathError(config InvalidPathConfig) InvalidPathError {
	return errors.Throw(errors.CtorConfig{
		Error: config.Err,
		Msg:   "invalid path",
		Meta: errors.MetaFields{
			"Path": config.Path,
		},
	})
}
