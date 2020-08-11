package metric

import (
	"context"
	"github.com/funnyecho/code-push/gateway/metric/interface/grpc/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func New(config *Config, fns ...func(*Options)) *client {
	clientOptions := &Options{}

	for _, fn := range fns {
		fn(clientOptions)
	}

	c := &client{
		Options: clientOptions,
	}

	return c
}

type client struct {
	*Options

	conn                  *grpc.ClientConn
	requestDurationClient pb.RequestDurationClient
}

func (c *client) RequestDuration(path string, success bool, durationSecond float64) {
	c.requestDurationClient.Gateway(context.Background(), &pb.GatewayRequestDurationRequest{
		Svr:            "client.g",
		Path:           path,
		Success:        success,
		DurationSecond: durationSecond,
	})
}

func (c *client) Conn() error {
	conn, err := grpc.Dial(
		c.Options.ServerAddr,
		grpc.WithInsecure(),
	)

	if err != nil {
		return errors.Wrapf(err, "Dail to grpc server: %c failed", c.Options.ServerAddr)
	}

	c.conn = conn
	c.requestDurationClient = pb.NewRequestDurationClient(conn)
	return nil
}

func (c *client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

type Config struct {
}

type Options struct {
	ServerAddr string
}
