package sessionAdapter

import (
	"github.com/funnyecho/code-push/daemon/session/interface/grpc/pb"
	"github.com/funnyecho/code-push/pkg/adapterkit"
	"github.com/funnyecho/code-push/pkg/adapterkit/grpc"
	"github.com/funnyecho/code-push/pkg/log"
	"google.golang.org/grpc"
)

func New(logger log.Logger, fns ...func(*Options)) *Client {
	ctorOptions := &Options{ServerAddr: ""}

	for _, fn := range fns {
		fn(ctorOptions)
	}

	var c *Client
	c = &Client{
		Adaptable: adapterkit_grpc.GrpcAdapter(
			adapterkit_grpc.WithGrpcAdaptName("session.d"),
			adapterkit_grpc.WithGrpcAdaptTarget(ctorOptions.ServerAddr),
			adapterkit_grpc.WithGrpcAdaptLogger(logger),
			adapterkit_grpc.WithGrpcAdaptConnected(func(conn *grpc.ClientConn) {
				c.accessTokenClient = pb.NewAccessTokenClient(conn)
			}),
		),
		Logger:  logger,
		Options: ctorOptions,
	}

	return c
}

type Client struct {
	log.Logger
	*Options
	adapterkit.Adaptable

	accessTokenClient pb.AccessTokenClient
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
