package sessionAdapter

import (
	"context"
	"github.com/funnyecho/code-push/daemon/session/interface/grpc/pb"
	"github.com/pkg/errors"
)

type AccessTokenIssuer = pb.AccessTokenIssuer

const (
	AccessTokenIssuer_SYS    = pb.AccessTokenIssuer_SYS
	AccessTokenIssuer_PORTAL = pb.AccessTokenIssuer_PORTAL
	AccessTokenIssuer_CLIENT = pb.AccessTokenIssuer_CLIENT
)

func (c *Client) GenerateAccessToken(issuer AccessTokenIssuer, subject string) ([]byte, error) {
	res, err := c.accessTokenClient.GenerateAccessToken(context.Background(), &pb.GenerateAccessTokenRequest{
		Claims: &pb.AccessTokenClaims{
			Issuer:   issuer,
			Subject:  subject,
			Audience: nil,
		},
	})

	return unmarshalStringResponse(res), err
}

func (c *Client) VerifyAccessToken(token string) (subject []byte, err error) {
	res, err := c.accessTokenClient.VerifyAccessToken(context.Background(), &pb.VerifyAccessTokenRequest{Token: token})

	if err != nil {
		return nil, err
	}

	claims := res.GetClaims()
	if claims == nil {
		return nil, errors.Wrap(ErrInvalidToken, "failed to fetch token claims")
	}

	return []byte(claims.Subject), nil
}
