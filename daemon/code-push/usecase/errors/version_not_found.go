package errors

import "github.com/funnyecho/code-push/pkg/errors"

type VersionNotFoundError error

func ThrowVersionNotFoundError(envId, appVersion string, err error) EnvNotFoundError {
	return errors.Throw(errors.CtorConfig{
		Error: err,
		Code:  FA_VERSION_NOT_FOUND,
		Msg:   "version not found",
		Meta:  errors.MetaFields{"envId": envId, "appVersion": appVersion},
	})
}
