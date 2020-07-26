package usecase

import "github.com/funnyecho/code-push/daemon/session"

type UseCase interface {
	AccessToken
}

type AccessToken interface {
	GenerateAccessToken(claims *session.AccessTokenClaims) ([]byte, error)
	VerifyAccessToken(token []byte) (*session.AccessTokenClaims, error)
}
