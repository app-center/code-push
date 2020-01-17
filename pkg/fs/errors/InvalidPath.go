package errors

import "github.com/funnyecho/code-push/pkg/errors"

type InvalidPathError errors.Error

type InvalidPathConfig struct {
	Err  error
	Path interface{}
}

func NewInvalidPathError(config InvalidPathConfig) *InvalidPathError {
	return &InvalidPathError{
		OpenError: errors.NewOpenError(errors.CtorConfig{
			Error: config.Err,
			Msg:   "invalid path",
			Meta: errors.MetaFields{
				"Path": config.Path,
			},
		}),
	}
}
