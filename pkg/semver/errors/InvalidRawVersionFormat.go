package errors

import (
	"github.com/funnyecho/code-push/pkg/errors"
)

type InvalidRawVersionFormatError error

type InvalidRawVersionFormatErrorConfig struct {
	RawVersion string
}

func NewInvalidRawVersionFormatError(config InvalidRawVersionFormatErrorConfig) InvalidRawVersionFormatError {
	return errors.Throw(errors.CtorConfig{
		Msg: "invalid raw version format",
		Meta: errors.MetaFields{
			"RawVersion": config.RawVersion,
		},
	})
}
