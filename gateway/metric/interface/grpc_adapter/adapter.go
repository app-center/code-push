package metricAdapter

import (
	"github.com/funnyecho/code-push/gateway/metric/interface/grpc/pb"
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
			adapterkit.WithGrpcAdaptName("metric.g"),
			adapterkit.WithGrpcAdaptTarget(ctorOptions.ServerAddr),
			adapterkit.WithGrpcAdaptLogger(logger),
			adapterkit.WithGrpcAdaptConnected(func(conn *grpc.ClientConn) {
				c.requestDurationClient = pb.NewRequestDurationClient(conn)
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

	requestDurationClient pb.RequestDurationClient
}

type Options struct {
	ServerAddr string
}
