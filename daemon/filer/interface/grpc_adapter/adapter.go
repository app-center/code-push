package filerAdapter

import (
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
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
			adapterkit.WithGrpcAdaptName("filer.d"),
			adapterkit.WithGrpcAdaptTarget(ctorOptions.ServerAddr),
			adapterkit.WithGrpcAdaptLogger(logger),
			adapterkit.WithGrpcAdaptConnected(func(conn *grpc.ClientConn) {
				c.uploadClient = pb.NewUploadClient(conn)
				c.fileClient = pb.NewFileClient(conn)
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

	uploadClient pb.UploadClient
	fileClient   pb.FileClient
}

type Options struct {
	ServerAddr string
}
