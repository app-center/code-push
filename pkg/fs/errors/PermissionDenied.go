package errors

import "github.com/funnyecho/code-push/pkg/errors"

type PermissionDeniedError errors.Error

type PermissionDeniedConfig struct {
	Err  error
	Path interface{}
}

func NewPermissionDeniedError(config PermissionDeniedConfig) *PermissionDeniedError {
	return &PermissionDeniedError{
		OpenError: errors.NewOpenError(errors.CtorConfig{
			Error: config.Err,
			Msg:   "permission denied",
			Meta: errors.MetaFields{
				"Path": config.Path,
			},
		}),
	}
}
