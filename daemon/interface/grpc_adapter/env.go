package adapter

import (
	"context"
	"github.com/funnyecho/code-push/daemon/interface/grpc/pb"
)

func (c *Client) CreateEnv(ctx context.Context, branchId, envId, envName, envEncToken []byte) (*pb.EnvResponse, error) {
	return c.envClient.CreateEnv(ctx, &pb.CreateEnvRequest{
		BranchId:    branchId,
		EnvName:     envName,
		EnvEncToken: envEncToken,
		EnvId:       envId,
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

func (c *Client) GetEnvsWithBranchId(ctx context.Context, branchId string) ([]*pb.EnvResponse, error) {
	res, err := c.envClient.GetEnvsWithBranchId(ctx, &pb.BranchIdRequest{BranchId: branchId})
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil
	}

	return res.List, nil
}
