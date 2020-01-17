package errors

import (
	"github.com/funnyecho/code-push/pkg/errors"
)

type InvalidPatchVersionError errors.Error

type InvalidPatchVersionErrorConfig struct {
	Err          error
	RawVersion   string
	PatchVersion interface{}
}

func NewInvalidPatchVersionError(config InvalidPatchVersionErrorConfig) *InvalidPatchVersionError {
	return &InvalidPatchVersionError{
		OpenError: errors.NewOpenError(errors.CtorConfig{
			Error: config.Err,
			Msg:   "invalid patch version",
			Meta: errors.MetaFields{
				"RawVersion":   config.RawVersion,
				"PatchVersion": config.PatchVersion,
			},
		}),
	}
}
