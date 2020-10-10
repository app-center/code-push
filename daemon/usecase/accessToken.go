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

	uc.accessTokenCache.Set(token, &accessTokenCacheEntry{
		claims:    claims,
		expiredAt: time.Now().Add(accessTokenExpiration).Unix(),
	}, cache.DefaultExpiration)

	return []byte(token), nil
}

func (uc *useCase) EvictAccessToken(token []byte) error {
	if token == nil {
		return errors.Wrap(daemon.ErrParamsInvalid, "token is required")
	}

	uc.accessTokenCache.Delete(string(token))
	return nil
}

func (uc *useCase) VerifyAccessToken(token []byte) (*daemon.AccessTokenClaims, error) {
	if token == nil {
		return nil, errors.Wrap(daemon.ErrParamsInvalid, "token is required")
	}

	v, existed := uc.accessTokenCache.Get(string(token))

	if !existed {
		return nil, daemon.ErrAccessTokenInvalid
	}

	entry := v.(*accessTokenCacheEntry)
	if entry == nil {
		return nil, daemon.ErrAccessTokenInvalid
	}

	claims := entry.claims
	expired := entry.expiredAt

	if claims == nil {
		return nil, daemon.ErrAccessTokenInvalid
	}

	if !uc.isValidAccessTokenIssuer(claims.Issuer) {
		return nil, errors.Wrapf(daemon.ErrAccessTokenInvalid, "invalid token issuer: %d", claims.Issuer)
	}

	now := time.Now().Unix()
	if now > expired {
		return nil, daemon.ErrAccessTokenInvalid
	} else {
		entry.expiredAt = time.Now().Add(accessTokenExpiration).Unix()
	}

	return claims, nil
}

func (uc *useCase) initAccessTokenUseCase() {
	uc.accessTokenCache = cache.New(accessTokenCacheExpiration, accessTokenCacheExpiration)
}

func (uc *useCase) isValidAccessTokenIssuer(issuer daemon.AccessTokenIssuer) bool {
	switch issuer {
	case daemon.AccessTokenIssuerSYS, daemon.AccessTokenIssuerPORTAL, daemon.AccessTokenIssuerCLIENT:
		return true
	default:
		return false
	}
}

const accessTokenCacheExpiration = 24 * time.Hour
const accessTokenExpiration = 1 * time.Hour

type accessTokenCacheEntry struct {
	claims    *daemon.AccessTokenClaims
	expiredAt int64
}
