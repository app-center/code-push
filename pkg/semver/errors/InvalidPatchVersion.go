package errors

import (
	"fmt"
	"github.com/funnyecho/code-push/pkg/errors"
)

type InvalidPatchVersionError struct {
	Err          error
	RawVersion   string
	PatchVersion string
}

func (err *InvalidPatchVersionError) Error() *errors.Error {
	return errors.New(errors.CtorConfig{
		Error: err.Err,
		Msg:   "invalid patch version",
		Meta: errors.MetaFields{
			"RawVersion":   err.RawVersion,
			"PatchVersion": err.PatchVersion,
		},
	})
}

func NewInvalidPatchVersionError(err error, rawVersion string, patchVersion interface{}) *InvalidPatchVersionError {
	return &InvalidPatchVersionError{
		Err:          err,
		RawVersion:   rawVersion,
		PatchVersion: fmt.Sprintf("%v", patchVersion),
	}
}
