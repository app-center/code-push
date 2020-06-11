package bolt

import "github.com/funnyecho/code-push/daemon/code-push/domain"

const (
	ErrBranchCreationParamsInvalid = domain.Error("invalid branch creation params")
	ErrBranchExists                = domain.Error("branch was existed")

	ErrBranchNotFound = domain.Error("branch not found")
)
