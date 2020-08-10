package session

import (
	"context"
	"github.com/funnyecho/code-push/daemon/session/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway/sys"
	"github.com/funnyecho/code-push/pkg/grpcInterceptor"
	"github.com/funnyecho/code-push/pkg/log"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
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

func (c *Client) GenerateAccessToken(subject string) ([]byte, error) {
	res, err := c.accessTokenClient.GenerateAccessToken(context.Background(), &pb.GenerateAccessTokenRequest{
		Claims: &pb.AccessTokenClaims{
			Issuer:   pb.AccessTokenIssuer_PORTAL,
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
		return nil, errors.Wrap(sys.ErrInvalidToken, "failed to fetch token claims")
	}

	return []byte(claims.Subject), nil
}

func (c *Client) Conn() error {
	conn, err := grpc.Dial(
		c.Options.ServerAddr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpcInterceptor.UnaryClientMetricInterceptor(c.Logger),
			grpcInterceptor.UnaryClientErrorInterceptor(),
		)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpcInterceptor.StreamClientMetricInterceptor(c.Logger),
			grpcInterceptor.StreamClientErrorInterceptor(),
		)),
	)
	if err != nil {
		return errors.Wrapf(err, "Dail to grpc server: %c failed", c.Options.ServerAddr)
	}

	c.conn = conn
	c.accessTokenClient = pb.NewAccessTokenClient(conn)
	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
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
