package usecase

import (
	"context"
	"crypto/md5"
	daemonAdapter "github.com/funnyecho/code-push/daemon/interface/grpc_adapter"
	"github.com/funnyecho/code-push/gateway"
	"github.com/funnyecho/code-push/pkg/jwt"
	"github.com/funnyecho/code-push/pkg/oauth"
	"github.com/funnyecho/code-push/pkg/util"
	"github.com/pkg/errors"
)

func (uc *useCase) AuthRootUser(_ context.Context, name, pwd string) error {
	if name == "" || pwd == "" {
		return errors.Wrap(gateway.ErrParamsInvalid, "name and pwd are required")
	}

	if name != uc.RootUserName || pwd != uc.RootUserPwd {
		return gateway.ErrUnauthorized
	}
	return nil
}

func (uc *useCase) SignTokenForRootUser(ctx context.Context) ([]byte, error) {
	salt := md5.Sum([]byte(uc.RootUserPwd))
	subject, subjectErr := util.EncryptAES(salt[:], []byte(uc.RootUserName))
	if subjectErr != nil {
		return nil, errors.WithStack(subjectErr)
	}

	token, tokenErr := uc.daemon.GenerateAccessToken(ctx, daemonAdapter.AccessTokenIssuerSys, string(subject))
	if tokenErr != nil {
		return nil, errors.WithStack(tokenErr)
	}

	return token, nil
}

func (uc *useCase) VerifyTokenForRootUser(ctx context.Context, token []byte) error {
	if token == nil {
		return gateway.ErrParamsInvalid
	}

	subject, verifyErr := uc.daemon.VerifyAccessToken(ctx, string(token))
	if verifyErr != nil {
		return errors.WithStack(verifyErr)
	}

	salt := md5.Sum([]byte(uc.RootUserPwd))
	plainSubject, plainSubjectErr := util.DecryptAES(salt[:], subject)
	if plainSubjectErr != nil || string(plainSubject) != uc.RootUserName {
		return errors.Wrap(gateway.ErrInvalidToken, "failed to verify subject in jwt token")
	}

	return nil
}

func (uc *useCase) AuthBranch(ctx context.Context, branchId, timestamp, nonce, sign []byte) error {
	if branchId == nil || timestamp == nil || nonce == nil || sign == nil {
		return errors.Wrap(gateway.ErrParamsInvalid, "branchId, timestamp, nonce or sign are required")
	}

	encToken, encTokenErr := uc.daemon.GetBranchEncToken(ctx, branchId)
	if encTokenErr != nil {
		return errors.Wrap(encTokenErr, "failed to get branch auth token")
	}

	authValid, authErr := oauth.Valid(string(encToken), string(timestamp), string(nonce), string(sign))
	if authErr != nil {
		return errors.WithStack(authErr)
	}

	if !authValid {
		return errors.WithStack(gateway.ErrUnauthorized)
	}

	return nil
}

func (uc *useCase) AuthBranchWithJWT(ctx context.Context, token string) (branchId []byte, err error) {
	claims, err := jwt.ExtractClaims(token)

	if err != nil {
		return nil, errors.WithMessagef(err, "failed to extract claims from token: %s", token)
	}

	iBranchId := []byte(claims.Subject)
	branchEncToken, branchEncTokenErr := uc.daemon.GetBranchEncToken(ctx, iBranchId)
	if branchEncTokenErr != nil {
		return nil, errors.WithMessagef(gateway.ErrUnauthorized, "failed to get token from branchId:%s", claims.Subject)
	}

	_, verifyErr := jwt.VerifyWithHMAC(token, branchEncToken)
	if verifyErr != nil {
		return nil, errors.WithMessage(gateway.ErrUnauthorized, "failed to verify token")
	}

	branchId = iBranchId
	err = nil

	return
}

func (uc *useCase) SignTokenForBranch(ctx context.Context, branchId []byte) ([]byte, error) {
	if branchId == nil {
		return nil, errors.Wrap(gateway.ErrParamsInvalid, "branchId required")
	}

	token, tokenErr := uc.daemon.GenerateAccessToken(ctx, daemonAdapter.AccessTokenIssuerPortal, string(branchId))
	if tokenErr != nil {
		return nil, errors.WithStack(tokenErr)
	}

	return token, nil
}

func (uc *useCase) VerifyTokenForBranch(ctx context.Context, token []byte) (branchId []byte, err error) {
	if token == nil {
		return nil, gateway.ErrParamsInvalid
	}

	subject, verifyErr := uc.daemon.VerifyAccessToken(ctx, string(token))
	if verifyErr != nil {
		return nil, errors.WithStack(verifyErr)
	}

	return subject, nil
}

func (uc *useCase) AuthEnv(ctx context.Context, envId, timestamp, nonce, sign []byte) error {
	if envId == nil || timestamp == nil || nonce == nil || sign == nil {
		return errors.Wrap(gateway.ErrParamsInvalid, "envId, timestamp, nonce or sign are required")
	}

	encToken, encTokenErr := uc.daemon.GetEnvEncToken(ctx, envId)
	if encTokenErr != nil {
		return errors.Wrap(encTokenErr, "failed to get env auth token")
	}

	authValid, authErr := oauth.Valid(string(encToken), string(timestamp), string(nonce), string(sign))
	if authErr != nil {
		return errors.WithStack(authErr)
	}

	if !authValid {
		return errors.WithStack(gateway.ErrUnauthorized)
	}

	return nil
}

func (uc *useCase) SignTokenForEnv(ctx context.Context, envId []byte) ([]byte, error) {
	if envId == nil {
		return nil, errors.Wrap(gateway.ErrParamsInvalid, "branchId required")
	}

	token, tokenErr := uc.daemon.GenerateAccessToken(ctx, daemonAdapter.AccessTokenIssuerClient, string(envId))
	if tokenErr != nil {
		return nil, errors.WithStack(tokenErr)
	}

	return token, nil
}

func (uc *useCase) VerifyTokenForEnv(ctx context.Context, token []byte) (envId []byte, err error) {
	if token == nil {
		return nil, gateway.ErrParamsInvalid
	}

	subject, verifyErr := uc.daemon.VerifyAccessToken(ctx, string(token))
	if verifyErr != nil {
		return nil, errors.WithStack(verifyErr)
	}

	return subject, nil
}

func (uc *useCase) EvictToken(ctx context.Context, token []byte) error {
	if token == nil {
		return gateway.ErrParamsInvalid
	}

	return uc.daemon.EvictAccessToken(ctx, string(token))
}
