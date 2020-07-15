package errors

import "github.com/pkg/errors"

type PermissionDeniedError error

type PermissionDeniedConfig struct {
	Err  error
	Path interface{}
}

func NewPermissionDeniedError(config PermissionDeniedConfig) PermissionDeniedError {
	return errors.WithMessagef(config.Err, "permission denied, path: %s", config.Path)
}
