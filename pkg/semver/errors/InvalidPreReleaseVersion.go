package errors

import (
	"github.com/funnyecho/code-push/pkg/errors"
)

type InvalidPreReleaseVersionError errors.Error

type InvalidPreReleaseVersionErrorConfig struct {
	Err        error
	RawVersion string
	RawPR      string
	PRStage    interface{}
	PRVersion  interface{}
	PRBuild    interface{}
}

func NewInvalidPreReleaseVersionError(config InvalidPreReleaseVersionErrorConfig) *InvalidPreReleaseVersionError {
	return &InvalidPreReleaseVersionError{
		OpenError: errors.NewOpenError(errors.CtorConfig{
			Error: config.Err,
			Msg:   "invalid pre release version",
			Meta: errors.MetaFields{
				"RawVersion": config.RawVersion,
				"RawPR":      config.RawPR,
				"PRStage":    config.PRStage,
				"PRVersion":  config.PRVersion,
				"PRBuild":    config.PRBuild,
			},
		}),
	}
}
