package domain

const (
	ErrParamsInvalid = Error("invalid operation params")

	ErrBranchCreationParamsInvalid = Error("invalid branch creation params")
	ErrBranchExists                = Error("branch was existed")
	ErrBranchNotFound              = Error("branch not found")

	ErrEnvCreationParamsInvalid = Error("invalid env creation params")
	ErrEnvExists                = Error("env was existed")
	ErrEnvNotFound              = Error("env not found")

	ErrVersionCreationParamsInvalid = Error("invalid version creation params")
	ErrVersionExisted               = Error("version was existed")
)

// Error represents a WTF error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
