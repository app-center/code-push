package errors

import "github.com/funnyecho/code-push/pkg/errors"

type BranchInvalidEncTokenError error

func ThrowBranchInvalidEncTokenError(err error) BranchInvalidEncTokenError {
	return errors.Throw(errors.CtorConfig{
		Error: err,
		Code:  FA_BRANCH_INVALID_ENC_TOKEN,
	})
}
