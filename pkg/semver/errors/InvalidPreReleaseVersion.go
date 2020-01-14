package errors

import "github.com/funnyecho/code-push/pkg/errors"

type InvalidPreReleaseVersionError struct {
	Err        error
	RawVersion string
	PRStage    string
	PRVersion  string
	PRBuild    string
}

func (err *InvalidPreReleaseVersionError) Error() *errors.Error {
	return errors.New(errors.CtorConfig{
		Error: err.Err,
		Msg:   "invalid pre release version",
		Meta: errors.MetaFields{
			"RawVersion": err.RawVersion,
			"PRStage":    err.PRStage,
			"PRVersion":  err.PRVersion,
			"PRBuild":    err.PRBuild,
		},
	})
}
