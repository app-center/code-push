package errors

import "github.com/funnyecho/code-push/pkg/errors"

type BranchSaveError error

func ThrowBranchSaveError(err error, branch interface{}) BranchSaveError {
	return errors.Throw(errors.CtorConfig{
		Error: err,
		Code:  FA_BRANCH_CAN_NOT_SAVE,
		Meta:  MetaFields{"branch": branch},
	})
}
