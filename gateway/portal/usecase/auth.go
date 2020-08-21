package usecase

import (
	"context"
	sessionAdapter "github.com/funnyecho/code-push/daemon/session/interface/grpc_adapter"
	"github.com/funnyecho/code-push/gateway/portal"
	"github.com/funnyecho/code-push/pkg/oauth"
	"github.com/pkg/errors"
)

func (u *useCase) Auth(ctx context.Context, branchId, timestamp, nonce, sign []byte) error {
	if branchId == nil || timestamp == nil || nonce == nil || sign == nil {
		return errors.Wrap(portal.ErrParamsInvalid, "branchId, timestamp, nonce or sign are required")
	}

	encToken, encTokenErr := u.codePush.GetBranchEncToken(ctx, branchId)
	if encTokenErr != nil {
		return errors.Wrap(encTokenErr, "failed to get branch auth token")
	}

	authValid, authErr := oauth.Valid(string(encToken), string(timestamp), string(nonce), string(sign))
	if authErr != nil {
		return errors.WithStack(authErr)
	}

	if !authValid {
		return errors.WithStack(portal.ErrUnauthorized)
	}

	return nil
}

func (u *useCase) SignToken(ctx context.Context, branchId []byte) ([]byte, error) {
	if branchId == nil {
		return nil, errors.Wrap(portal.ErrParamsInvalid, "branchId required")
	}

	token, tokenErr := u.session.GenerateAccessToken(ctx, sessionAdapter.AccessTokenIssuerPortal, string(branchId))
	if tokenErr != nil {
		return nil, errors.WithStack(tokenErr)
	}

	return token, nil
}

func (u *useCase) VerifyToken(ctx context.Context, token []byte) (branchId []byte, err error) {
	if token == nil {
		return nil, portal.ErrParamsInvalid
	}

	subject, verifyErr := u.session.VerifyAccessToken(ctx, string(token))
	if verifyErr != nil {
		return nil, errors.WithStack(verifyErr)
	}

	return subject, nil
}
