package errors

import "github.com/funnyecho/code-push/pkg/errors"

type BranchDeleteFailedError error

func ThrowBranchDeleteFailedError(branchId string, reason string) BranchDeleteFailedError {
	return errors.Throw(errors.CtorConfig{
		Code: FA_BRANCH_DELETE_FAILED,
		Msg:  "failed to delete branch",
		Meta: errors.MetaFields{"branchId": branchId, "reason": reason},
	})
}
