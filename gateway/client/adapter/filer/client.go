package filer

import (
	"context"
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway/client"
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

	conn       *grpc.ClientConn
	fileClient pb.FileClient
}

func (s *Client) GetSource(fileKey []byte) ([]byte, error) {
	if fileKey == nil {
		return nil, client.ErrParamsInvalid
	}

	res, err := s.fileClient.GetSource(context.Background(), &pb.GetSourceRequest{Key: fileKey})

	return unmarshalStringResponse(res), err
}

func (s *Client) Conn() error {
	conn, err := grpc.Dial(s.Options.ServerAddr, grpc.WithInsecure())
	if err != nil {
		return errors.Wrapf(err, "Dail to grpc server: %s failed", s.Options.ServerAddr)
	}

	s.conn = conn
	s.fileClient = pb.NewFileClient(conn)
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
