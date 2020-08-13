package codePushAdapter

import (
	"context"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
)

func (c *Client) CreateEnv(branchId, envName []byte) (*pb.EnvResponse, error) {
	return c.envClient.CreateEnv(context.Background(), &pb.CreateEnvRequest{
		BranchId: branchId,
		EnvName:  envName,
	})
}

func (c *Client) GetEnv(envId []byte) (*pb.EnvResponse, error) {
	return c.envClient.GetEnv(context.Background(), &pb.EnvIdRequest{EnvId: envId})
}

func (c *Client) DeleteEnv(envId []byte) error {
	_, err := c.envClient.DeleteEnv(context.Background(), &pb.EnvIdRequest{EnvId: envId})
	return err
}

func (c *Client) GetEnvEncToken(envId []byte) ([]byte, error) {
	res, err := c.envClient.GetEnvEncToken(context.Background(), &pb.EnvIdRequest{EnvId: envId})
	return unmarshalStringResponse(res), err
}
