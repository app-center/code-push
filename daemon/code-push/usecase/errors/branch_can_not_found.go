package errors

import "github.com/funnyecho/code-push/pkg/errors"

type BranchCanNotFoundError error

func ThrowBranchCanNotFoundError(branchId string) BranchCanNotFoundError {
	return errors.Throw(errors.CtorConfig{
		Code: FA_BRANCH_CAN_NOT_FOUND,
		Msg:  "branch name existed",
		Meta: errors.MetaFields{"branchId": branchId},
	})
}
