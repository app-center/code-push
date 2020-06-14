package code_push

const (
	ErrParamsInvalid = Error("FA_PARAMS_INVALID")

	ErrBranchExisted  = Error("FA_BRANCH_EXISTED")
	ErrBranchNotFound = Error("FA_BRANCH_NOT_FOUND")

	ErrEnvExisted  = Error("FA_ENV_EXISTED")
	ErrEnvNotFound = Error("FA_ENV_NOT_FOUND")

	ErrVersionExisted = Error("FA_VERSION_EXISTED")
)

// Error represents a WTF error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
