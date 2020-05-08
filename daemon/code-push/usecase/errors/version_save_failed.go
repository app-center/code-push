package errors

import "github.com/funnyecho/code-push/pkg/errors"

type VersionSaveError error

func ThrowVersionSaveError(err error, msg string, meta MetaFields) EnvSaveError {
	return errors.Throw(errors.CtorConfig{
		Error: err,
		Code:  FA_VERSION_SAVE_FAILED,
		Msg:   msg,
		Meta:  meta,
	})
}
