package errors

import "github.com/funnyecho/code-push/pkg/errors"

type BranchSaveError error

func ThrowBranchSaveError(err error, branch interface{}) BranchSaveError {
	return errors.Throw(errors.CtorConfig{
		Error: err,
		Code:  FA_BRANCH_SAVE_FAILED,
		Meta:  MetaFields{"branch": branch},
	})
}
