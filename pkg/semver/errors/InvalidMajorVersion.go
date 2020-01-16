package errors

import (
	"fmt"
	"github.com/funnyecho/code-push/pkg/errors"
)

type InvalidMajorVersionError struct {
	Err          error
	RawVersion   string
	MajorVersion string
}

func (err *InvalidMajorVersionError) Error() *errors.Error {
	return errors.New(errors.CtorConfig{
		Error: err.Err,
		Msg:   "invalid major version",
		Meta: errors.MetaFields{
			"RawVersion":   err.RawVersion,
			"MajorVersion": err.MajorVersion,
		},
	})
}

func NewInvalidMajorVersionError(err error, rawVersion string, majorVersion interface{}) *InvalidMajorVersionError {
	return &InvalidMajorVersionError{
		Err:          err,
		RawVersion:   rawVersion,
		MajorVersion: fmt.Sprintf("%v", majorVersion),
	}
}
