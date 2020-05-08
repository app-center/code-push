package errors

import "github.com/funnyecho/code-push/pkg/errors"

type VersionReleaseFailedError error

func ThrowVersionReleaseFailedError(err error, code string, params interface{}) VersionReleaseFailedError {
	if len(code) == 0 {
		code = FA_VERSION_RELEASE_FAILED
	}

	return errors.Throw(errors.CtorConfig{
		Error: err,
		Code:  code,
		Msg:   "version release failed",
		Meta:  errors.MetaFields{"params": params},
	})
}
