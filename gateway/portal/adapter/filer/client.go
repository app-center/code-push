package filer

import (
	"context"
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
	"github.com/funnyecho/code-push/pkg/grpcStreamer"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"io"
)

func New(fns ...func(*Options)) *Client {
	ctorOptions := &Options{ServerAddr: ""}

	for _, fn := range fns {
		fn(ctorOptions)
	}

	return &Client{
		options: ctorOptions,
	}
}

type Client struct {
	options *Options

	conn         *grpc.ClientConn
	uploadClient pb.UploadClient
}

func (s *Client) UploadPkg(source io.Reader) (fileKey []byte, err error) {
	stream, err := s.uploadClient.UploadToAliOss(context.Background())

	streamSender := grpcStreamer.NewSender(func(p byte) error {
		sendErr := stream.Send(&pb.UploadToAliOssRequest{Data: uint32(p)})
		if sendErr == io.EOF {
			return sendErr
		}

		return nil
	})

	written, copyErr := io.Copy(streamSender, source)
	if copyErr != nil {
		// FIXME: maybe stream need to close
		return nil, errors.Wrapf(err, "failed to write to client stream, written: %d", written)
	}

	res, resErr := stream.CloseAndRecv()
	return unmarshalStringResponse(res), resErr
}

func (s *Client) Conn() error {
	conn, err := grpc.Dial(s.options.ServerAddr, grpc.WithInsecure())
	if err != nil {
		return errors.Wrapf(err, "Dail to grpc server: %s failed", s.options.ServerAddr)
	}

	s.conn = conn
	s.uploadClient = pb.NewUploadClient(conn)
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
