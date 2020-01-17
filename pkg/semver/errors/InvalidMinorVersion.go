package errors

import (
	"github.com/funnyecho/code-push/pkg/errors"
)

type InvalidMinorVersionError errors.Error

type InvalidMinorVersionErrorConfig struct {
	Err          error
	RawVersion   string
	MinorVersion interface{}
}

func NewInvalidMinorVersionError(config InvalidMinorVersionErrorConfig) *InvalidMinorVersionError {
	return &InvalidMinorVersionError{
		OpenError: errors.NewOpenError(errors.CtorConfig{
			Error: config.Err,
			Msg:   "invalid minor version",
			Meta: errors.MetaFields{
				"RawVersion":   config.RawVersion,
				"MinorVersion": config.MinorVersion,
			},
		}),
	}
}
