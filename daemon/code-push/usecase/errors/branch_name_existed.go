package errors

import "github.com/funnyecho/code-push/pkg/errors"

type BranchNameExistedError error

func NewBranchNameExistedError(branchName string) BranchNameExistedError {
	return errors.Throw(errors.CtorConfig{
		Code: FA_BRANCH_NAME_EXISTED,
		Msg:  "branch name existed",
		Meta: errors.MetaFields{"branchName": branchName},
	})
}
