package errors

import (
	"github.com/funnyecho/code-push/pkg/errors"
)

type InvalidMajorVersionError struct {
	RawVersion   string
	MajorVersion string
}

func (err *InvalidMajorVersionError) Error() *errors.Error {
	return errors.New(errors.CtorConfig{
		Msg: "invalid major version",
		Meta: errors.MetaFields{
			"RawVersion":   err.RawVersion,
			"MajorVersion": err.MajorVersion,
		},
	})
}
