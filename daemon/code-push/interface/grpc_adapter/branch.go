package codePushAdapter

import (
	"context"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
)

func (c *Client) CreateBranch(ctx context.Context, branchName []byte) (*pb.BranchResponse, error) {
	return c.branchClient.CreateBranch(ctx, &pb.CreateBranchRequest{
		BranchName: branchName,
	})
}

func (c *Client) DeleteBranch(ctx context.Context, branchId []byte) error {
	_, err := c.branchClient.DeleteBranch(ctx, &pb.DeleteBranchRequest{BranchId: branchId})
	return err
}

func (c *Client) GetBranchEncToken(ctx context.Context, branchId []byte) ([]byte, error) {
	res, err := c.branchClient.GetBranchEncToken(ctx, &pb.GetBranchEncTokenRequest{BranchId: branchId})
	return unmarshalStringResponse(res), err
}
