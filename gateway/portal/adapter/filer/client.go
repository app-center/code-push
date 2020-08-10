package filer

import (
	"context"
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
	"github.com/funnyecho/code-push/pkg/grpcInterceptor"
	"github.com/funnyecho/code-push/pkg/grpcStreamer"
	"github.com/funnyecho/code-push/pkg/log"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"io"
	"mime/multipart"
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

	conn         *grpc.ClientConn
	uploadClient pb.UploadClient
}

func (c *Client) UploadPkg(source multipart.File) (fileKey []byte, err error) {
	stream, err := c.uploadClient.UploadToAliOss(context.Background())

	streamSender := grpcStreamer.NewSender(func(p byte) (err error) {
		err = stream.Send(&pb.UploadToAliOssRequest{Data: uint32(p)})

		return
	})

	written, copyErr := io.Copy(streamSender, source)
	if copyErr != nil {
		_ = stream.CloseSend()
		return nil, errors.Wrapf(copyErr, "failed to write to client stream, written: %d", written)
	}

	res, resErr := stream.CloseAndRecv()
	return unmarshalStringResponse(res), resErr
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
	c.uploadClient = pb.NewUploadClient(conn)
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
