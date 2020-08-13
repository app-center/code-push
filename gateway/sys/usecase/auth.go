package usecase

import (
	"crypto/md5"
	sessionAdapter "github.com/funnyecho/code-push/daemon/session/interface/grpc_adapter"
	"github.com/funnyecho/code-push/gateway/sys"
	"github.com/funnyecho/code-push/pkg/util"
	"github.com/pkg/errors"
)

func (u *useCase) Auth(name, pwd string) error {
	if name == "" || pwd == "" {
		return errors.Wrap(sys.ErrParamsInvalid, "name and pwd are required")
	}

	if name != u.options.RootUserName || pwd != u.options.RootUserPwd {
		return sys.ErrUnauthorized
	}
	return nil
}

func (u *useCase) SignToken() ([]byte, error) {
	salt := md5.Sum([]byte(u.options.RootUserPwd))
	subject, subjectErr := util.EncryptAES(salt[:], []byte(u.options.RootUserName))
	if subjectErr != nil {
		return nil, errors.WithStack(subjectErr)
	}

	token, tokenErr := u.session.GenerateAccessToken(sessionAdapter.AccessTokenIssuer_SYS, string(subject))
	if tokenErr != nil {
		return nil, errors.WithStack(tokenErr)
	}

	return token, nil
}

func (u *useCase) VerifyToken(token []byte) error {
	if token == nil {
		return sys.ErrParamsInvalid
	}

	subject, verifyErr := u.session.VerifyAccessToken(string(token))
	if verifyErr != nil {
		return errors.WithStack(verifyErr)
	}

	salt := md5.Sum([]byte(u.options.RootUserPwd))
	plainSubject, plainSubjectErr := util.DecryptAES(salt[:], subject)
	if plainSubjectErr != nil || string(plainSubject) != u.options.RootUserName {
		return errors.Wrap(sys.ErrInvalidToken, "failed to verify subject in jwt token")
	}

	return nil
}
