package errors

import "github.com/funnyecho/code-push/pkg/errors"

type InvalidMajorVersionError errors.Error

type InvalidMajorVersionErrorConfig struct {
	Err          error
	RawVersion   string
	MajorVersion interface{}
}

func NewInvalidMajorVersionError(config InvalidMajorVersionErrorConfig) *InvalidMajorVersionError {
	return &InvalidMajorVersionError{
		OpenError: errors.NewOpenError(errors.CtorConfig{
			Error: config.Err,
			Msg:   "invalid major version",
			Meta: errors.MetaFields{
				"RawVersion":   config.RawVersion,
				"MajorVersion": config.MajorVersion,
			},
		}),
	}
}
