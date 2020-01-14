package errors

import (
	"github.com/funnyecho/code-push/pkg/errors"
)

type InvalidMinorVersionError struct {
	Err          error
	RawVersion   string
	MinorVersion string
}

func (err *InvalidMinorVersionError) Error() *errors.Error {
	return errors.New(errors.CtorConfig{
		Error: err.Err,
		Msg:   "invalid minor version",
		Meta: errors.MetaFields{
			"RawVersion":   err.RawVersion,
			"MinorVersion": err.MinorVersion,
		},
	})
}
