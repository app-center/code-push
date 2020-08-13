package filerAdapter

import (
	"context"
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
)

func (c *Client) GetSource(fileKey []byte) ([]byte, error) {
	if fileKey == nil {
		return nil, ErrParamsInvalid
	}

	res, err := c.fileClient.GetSource(context.Background(), &pb.GetSourceRequest{Key: fileKey})

	return unmarshalStringResponse(res), err
}
