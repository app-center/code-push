package errors

import (
	"github.com/funnyecho/code-push/pkg/errors"
)

type InvalidRawVersionFormatError errors.Error

type InvalidRawVersionFormatErrorConfig struct {
	RawVersion string
}

func NewInvalidRawVersionFormatError(config InvalidRawVersionFormatErrorConfig) *InvalidRawVersionFormatError {
	return &InvalidRawVersionFormatError{
		OpenError: errors.NewOpenError(errors.CtorConfig{
			Msg: "invalid raw version format",
			Meta: errors.MetaFields{
				"RawVersion": config.RawVersion,
			},
		})}
}
