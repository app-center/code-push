package adapter

import (
	"context"
	"github.com/funnyecho/code-push/daemon/interface/grpc/pb"
	"github.com/pkg/errors"
)

type AccessTokenIssuer = pb.AccessTokenIssuer

const (
	AccessTokenIssuerSys    = pb.AccessTokenIssuer_SYS
	AccessTokenIssuerPortal = pb.AccessTokenIssuer_PORTAL
	AccessTokenIssuerClient = pb.AccessTokenIssuer_CLIENT
)

func (c *Client) GenerateAccessToken(ctx context.Context, issuer AccessTokenIssuer, subject string) ([]byte, error) {
	res, err := c.accessTokenClient.GenerateAccessToken(ctx, &pb.GenerateAccessTokenRequest{
		Claims: &pb.AccessTokenClaims{
			Issuer:   issuer,
			Subject:  subject,
			Audience: nil,
		},
	})

	return unmarshalStringResponse(res), err
}

func (c *Client) VerifyAccessToken(ctx context.Context, token string) (subject []byte, err error) {
	res, err := c.accessTokenClient.VerifyAccessToken(ctx, &pb.VerifyAccessTokenRequest{Token: token})

	if err != nil {
		return nil, err
	}

	claims := res.GetClaims()
	if claims == nil {
		return nil, errors.Wrap(ErrInvalidToken, "failed to fetch token claims")
	}

	return []byte(claims.Subject), nil
}

func (c *Client) EvictAccessToken(ctx context.Context, token string) error {
	_, err := c.accessTokenClient.EvictAccessToken(ctx, &pb.EvictAccessTokenRequest{Token: token})

	return err
}
