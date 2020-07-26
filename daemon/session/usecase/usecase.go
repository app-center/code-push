package usecase

import (
	"github.com/funnyecho/code-push/daemon/session"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

func New() UseCase {
	uc := &useCase{}

	uc.initCache()

	return uc
}

type useCase struct {
	accessTokenCache *cache.Cache
}

func (uc *useCase) GenerateAccessToken(claims *session.AccessTokenClaims) ([]byte, error) {
	if claims == nil {
		return nil, errors.Wrap(session.ErrParamsInvalid, "claims is required")
	}

	if !uc.isValidIssuer(claims.Issuer) {
		return nil, errors.Wrapf(session.ErrParamsInvalid, "invalid token issuer: %d", claims.Issuer)
	}

	if len(claims.Subject) == 0 {
		return nil, errors.Wrap(session.ErrParamsInvalid, "issuer or subject are required")
	}

	token := uuid.NewV4().String()

	uc.accessTokenCache.Set(token, claims, cache.DefaultExpiration)

	return []byte(token), nil
}

func (uc *useCase) VerifyAccessToken(token []byte) (*session.AccessTokenClaims, error) {
	if token == nil {
		return nil, errors.Wrap(session.ErrParamsInvalid, "claims is required")
	}

	v, existed := uc.accessTokenCache.Get(string(token))

	if !existed {
		return nil, session.ErrAccessTokenInvalid
	}

	claims := v.(*session.AccessTokenClaims)
	if claims == nil {
		return nil, session.ErrAccessTokenInvalid
	}

	if !uc.isValidIssuer(claims.Issuer) {
		return nil, errors.Wrapf(session.ErrAccessTokenInvalid, "invalid token issuer: %d", claims.Issuer)
	}

	return claims, nil
}

func (uc *useCase) initCache() {
	uc.accessTokenCache = cache.New(24*time.Hour, 24*time.Hour)
}

func (uc *useCase) isValidIssuer(issuer session.AccessTokenIssuer) bool {
	switch issuer {
	case session.AccessTokenIssuerSYS, session.AccessTokenIssuerPORTAL, session.AccessTokenIssuerCLIENT:
		return true
	default:
		return false
	}
}
