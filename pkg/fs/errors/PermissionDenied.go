package errors

import "github.com/funnyecho/code-push/pkg/errors"

type PermissionDeniedError error

type PermissionDeniedConfig struct {
	Err  error
	Path interface{}
}

func NewPermissionDeniedError(config PermissionDeniedConfig) PermissionDeniedError {
	return errors.Throw(errors.CtorConfig{
		Error: config.Err,
		Msg:   "permission denied",
		Meta: errors.MetaFields{
			"Path": config.Path,
		},
	})
}
