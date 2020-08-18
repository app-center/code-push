package filerAdapter

import (
	"context"
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
)

func (c *Client) GetSource(fileKey []byte) (*pb.FileSource, error) {
	if fileKey == nil {
		return nil, ErrParamsInvalid
	}

	return c.fileClient.GetSource(context.Background(), &pb.GetSourceRequest{Key: string(fileKey)})
}
