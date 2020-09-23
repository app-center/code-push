package usecase

import (
	"github.com/funnyecho/code-push/daemon"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

func (uc *useCase) GenerateAccessToken(claims *daemon.AccessTokenClaims) ([]byte, error) {
	if claims == nil {
		return nil, errors.Wrap(daemon.ErrParamsInvalid, "claims is required")
	}

	if !uc.isValidAccessTokenIssuer(claims.Issuer) {
		return nil, errors.Wrapf(daemon.ErrParamsInvalid, "invalid token issuer: %d", claims.Issuer)
	}

	if len(claims.Subject) == 0 {
		return nil, errors.Wrap(daemon.ErrParamsInvalid, "issuer or subject are required")
	}

	token := uuid.NewV4().String()

	uc.accessTokenCache.Set(token, claims, cache.DefaultExpiration)

	return []byte(token), nil
}

func (uc *useCase) VerifyAccessToken(token []byte) (*daemon.AccessTokenClaims, error) {
	if token == nil {
		return nil, errors.Wrap(daemon.ErrParamsInvalid, "claims is required")
	}

	v, existed := uc.accessTokenCache.Get(string(token))

	if !existed {
		return nil, daemon.ErrAccessTokenInvalid
	}

	claims := v.(*daemon.AccessTokenClaims)
	if claims == nil {
		return nil, daemon.ErrAccessTokenInvalid
	}

	if !uc.isValidAccessTokenIssuer(claims.Issuer) {
		return nil, errors.Wrapf(daemon.ErrAccessTokenInvalid, "invalid token issuer: %d", claims.Issuer)
	}

	return claims, nil
}

func (uc *useCase) initAccessTokenUseCase() {
	uc.accessTokenCache = cache.New(24*time.Hour, 24*time.Hour)
}

func (uc *useCase) isValidAccessTokenIssuer(issuer daemon.AccessTokenIssuer) bool {
	switch issuer {
	case daemon.AccessTokenIssuerSYS, daemon.AccessTokenIssuerPORTAL, daemon.AccessTokenIssuerCLIENT:
		return true
	default:
		return false
	}
}
