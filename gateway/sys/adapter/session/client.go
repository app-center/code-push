package session

import (
	"context"
	"github.com/funnyecho/code-push/daemon/session/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway/portal"
	"github.com/funnyecho/code-push/pkg/log"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func New(logger log.Logger, fns ...func(*Options)) *Client {
	ctorOptions := &Options{ServerAddr: ""}

	for _, fn := range fns {
		fn(ctorOptions)
	}

	return &Client{
		Logger:  logger,
		Options: ctorOptions,
	}
}

type Client struct {
	log.Logger
	*Options

	conn              *grpc.ClientConn
	accessTokenClient pb.AccessTokenClient
}

func (s *Client) GenerateAccessToken(subject string) ([]byte, error) {
	res, err := s.accessTokenClient.GenerateAccessToken(context.Background(), &pb.GenerateAccessTokenRequest{
		Claims: &pb.AccessTokenClaims{
			Issuer:   pb.AccessTokenIssuer_SYS,
			Subject:  subject,
			Audience: nil,
		},
	})

	return unmarshalStringResponse(res), err
}

func (s *Client) VerifyAccessToken(token string) (subject []byte, err error) {
	res, err := s.accessTokenClient.VerifyAccessToken(context.Background(), &pb.VerifyAccessTokenRequest{Token: token})

	if err != nil {
		return nil, err
	}

	claims := res.GetClaims()
	if claims == nil {
		return nil, errors.Wrap(portal.ErrInvalidToken, "failed to fetch token claims")
	}

	return []byte(claims.Subject), nil
}

func (s *Client) Conn() error {
	conn, err := grpc.Dial(s.Options.ServerAddr, grpc.WithInsecure())
	if err != nil {
		return errors.Wrapf(err, "Dail to grpc server: %s failed", s.Options.ServerAddr)
	}

	s.conn = conn
	s.accessTokenClient = pb.NewAccessTokenClient(conn)
	return nil
}

func (s *Client) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}

	return nil
}

type Options struct {
	ServerAddr string
}

func unmarshalStringResponse(r *pb.StringResponse) []byte {
	if r == nil {
		return nil
	}

	return []byte(r.Data)
}
