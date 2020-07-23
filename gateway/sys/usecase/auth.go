package usecase

import (
	"github.com/funnyecho/code-push/gateway/sys"
	"github.com/funnyecho/code-push/pkg/util"
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

func (u *useCase) SignToken() ([]byte, error) {
	subject, subjectErr := util.EncryptAES([]byte(u.options.RootUserPwd), []byte(u.options.RootUserName))
	if subjectErr != nil {
		return nil, errors.WithStack(subjectErr)
	}

	token, tokenErr := u.jwt.SignToken(string(subject))
	if tokenErr != nil {
		return nil, errors.WithStack(tokenErr)
	}

	return []byte(token), nil
}

func (u *useCase) VerifyToken(token []byte) error {
	if token == nil {
		return sys.ErrParamsInvalid
	}

	claims, verifyErr := u.jwt.VerifyToken(string(token))
	if verifyErr != nil {
		return errors.WithStack(verifyErr)
	}

	plainSubject, plainSubjectErr := util.DecryptAES([]byte(u.options.RootUserPwd), []byte(claims.Subject))
	if plainSubjectErr != nil || string(plainSubject) != u.options.RootUserName {
		return errors.Wrap(sys.ErrInvalidToken, "failed to verify subject in jwt token")
	}

	return nil
}
