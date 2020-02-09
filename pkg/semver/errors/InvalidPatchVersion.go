package errors

import (
	"github.com/funnyecho/code-push/pkg/errors"
)

type InvalidPatchVersionError error

type InvalidPatchVersionErrorConfig struct {
	Err          error
	RawVersion   string
	PatchVersion interface{}
}

func NewInvalidPatchVersionError(config InvalidPatchVersionErrorConfig) InvalidPatchVersionError {
	return errors.Throw(errors.CtorConfig{
		Error: config.Err,
		Msg:   "invalid patch version",
		Meta: errors.MetaFields{
			"RawVersion":   config.RawVersion,
			"PatchVersion": config.PatchVersion,
		},
	})
}
