package usecase

import (
	"context"
	sessionAdapter "github.com/funnyecho/code-push/daemon/session/interface/grpc_adapter"
	"github.com/funnyecho/code-push/gateway/client"
	"github.com/funnyecho/code-push/pkg/oauth"
	"github.com/pkg/errors"
)

func (uc *useCase) Auth(ctx context.Context, envId, timestamp, nonce, sign []byte) error {
	if envId == nil || timestamp == nil || nonce == nil || sign == nil {
		return errors.Wrap(client.ErrParamsInvalid, "envId, timestamp, nonce or sign are required")
	}

	encToken, encTokenErr := uc.codePush.GetEnvEncToken(ctx, envId)
	if encTokenErr != nil {
		return errors.Wrap(encTokenErr, "failed to get env auth token")
	}

	authValid, authErr := oauth.Valid(string(encToken), string(timestamp), string(nonce), string(sign))
	if authErr != nil {
		return errors.WithStack(authErr)
	}

	if !authValid {
		return errors.WithStack(client.ErrUnauthorized)
	}

	return nil
}

func (uc *useCase) SignToken(ctx context.Context, envId []byte) ([]byte, error) {
	if envId == nil {
		return nil, errors.Wrap(client.ErrParamsInvalid, "branchId required")
	}

	token, tokenErr := uc.session.GenerateAccessToken(ctx, sessionAdapter.AccessTokenIssuerClient, string(envId))
	if tokenErr != nil {
		return nil, errors.WithStack(tokenErr)
	}

	return token, nil
}

func (uc *useCase) VerifyToken(ctx context.Context, token []byte) (envId []byte, err error) {
	if token == nil {
		return nil, client.ErrParamsInvalid
	}

	subject, verifyErr := uc.session.VerifyAccessToken(ctx, string(token))
	if verifyErr != nil {
		return nil, errors.WithStack(verifyErr)
	}

	return subject, nil
}
