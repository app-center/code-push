package errors

import "github.com/funnyecho/code-push/pkg/errors"

const (
	FA_INVALID_PARAMS = "FA_INVALID_PARAMS"

	FA_BRANCH_NAME_EXISTED      = "FA_BRANCH_NAME_EXISTED"
	FA_BRANCH_INVALID_ENC_TOKEN = "FA_BRANCH_INVALID_ENC_TOKEN"
	FA_BRANCH_CAN_NOT_SAVE      = "FA_BRANCH_CAN_NOT_SAVE"
)

type MetaFields = errors.MetaFields
