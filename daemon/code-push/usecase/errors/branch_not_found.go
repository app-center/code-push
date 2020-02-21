package errors

import "github.com/funnyecho/code-push/pkg/errors"

type BranchNotFoundError error

func ThrowBranchNotFoundError(branchId string, err error) BranchNotFoundError {
	return errors.Throw(errors.CtorConfig{
		Error: err,
		Code:  FA_BRANCH_CAN_NOT_FOUND,
		Msg:   "branch not found",
		Meta:  errors.MetaFields{"branchId": branchId},
	})
}
