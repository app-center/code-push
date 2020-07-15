package errors

import "github.com/pkg/errors"

type InvalidPathError error

type InvalidPathConfig struct {
	Err  error
	Path interface{}
}

func NewInvalidPathError(config InvalidPathConfig) InvalidPathError {
	return errors.WithMessagef(config.Err, "invalid path: %s", config.Path)
}
