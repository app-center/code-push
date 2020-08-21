package codePushAdapter

import (
	"context"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
)

func (c *Client) CreateEnv(ctx context.Context, branchId, envName []byte) (*pb.EnvResponse, error) {
	return c.envClient.CreateEnv(ctx, &pb.CreateEnvRequest{
		BranchId: branchId,
		EnvName:  envName,
	})
}

func (c *Client) GetEnv(ctx context.Context, envId []byte) (*pb.EnvResponse, error) {
	return c.envClient.GetEnv(ctx, &pb.EnvIdRequest{EnvId: envId})
}

func (c *Client) DeleteEnv(ctx context.Context, envId []byte) error {
	_, err := c.envClient.DeleteEnv(ctx, &pb.EnvIdRequest{EnvId: envId})
	return err
}

func (c *Client) GetEnvEncToken(ctx context.Context, envId []byte) ([]byte, error) {
	res, err := c.envClient.GetEnvEncToken(ctx, &pb.EnvIdRequest{EnvId: envId})
	return unmarshalStringResponse(res), err
}
