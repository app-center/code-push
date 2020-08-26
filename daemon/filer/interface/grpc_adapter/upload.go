package filerAdapter

import (
	"context"
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
	"github.com/funnyecho/code-push/pkg/grpc-streamer"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
)

func (c *Client) UploadPkg(ctx context.Context, source multipart.File) (fileKey []byte, err error) {
	stream, err := c.uploadClient.UploadToAliOss(ctx)

	streamSender := grpc_streamer.NewSender(func(p byte) (err error) {
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
