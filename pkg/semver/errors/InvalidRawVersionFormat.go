package errors

import (
	"github.com/funnyecho/code-push/pkg/errors"
)

type InvalidRawVersionFormatError struct {
	RawVersion string
}

func (err *InvalidRawVersionFormatError) Error() *errors.Error {
	return errors.New(errors.CtorConfig{
		Msg: "invalid raw version format",
		Meta: errors.MetaFields{
			"RawVersion": err.RawVersion,
		},
	})
}
