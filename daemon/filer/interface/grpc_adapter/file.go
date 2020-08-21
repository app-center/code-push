package filerAdapter

import (
	"context"
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
)

func (c *Client) GetSource(ctx context.Context, fileKey []byte) (*pb.FileSource, error) {
	if fileKey == nil {
		return nil, ErrParamsInvalid
	}

	return c.fileClient.GetSource(ctx, &pb.GetSourceRequest{Key: string(fileKey)})
}
