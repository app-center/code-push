package usecase

import (
	"github.com/funnyecho/code-push/gateway/sys"
	"github.com/pkg/errors"
)

func (u *useCase) Auth(name, pwd []byte) error {
	if name == nil || pwd == nil {
		return errors.Wrap(sys.ErrParamsInvalid, "name and pwd are required")
	}

	if string(name) != u.options.RootUserName || string(pwd) != u.options.RootUserPwd {
		return sys.ErrUnauthorized
	}

	return nil
}
