package filerAdapter

import (
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	"github.com/funnyecho/code-push/pkg/adapterkit"
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
		Adaptable: adapterkit.GrpcAdapter(
			adapterkit.WithGrpcAdaptName("code-push.d"),
			adapterkit.WithGrpcAdaptTarget(ctorOptions.ServerAddr),
			adapterkit.WithGrpcAdaptLogger(logger),
			adapterkit.WithGrpcAdaptConnected(func(conn *grpc.ClientConn) {
				c.branchClient = pb.NewBranchClient(conn)
				c.envClient = pb.NewEnvClient(conn)
				c.versionClient = pb.NewVersionClient(conn)
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

	branchClient  pb.BranchClient
	envClient     pb.EnvClient
	versionClient pb.VersionClient
}

type Options struct {
	ServerAddr string
}
