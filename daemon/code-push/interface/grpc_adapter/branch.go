package codePushAdapter

import (
	"context"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
)

func (c *Client) CreateBranch(branchName []byte) (*pb.BranchResponse, error) {
	return c.branchClient.CreateBranch(context.Background(), &pb.CreateBranchRequest{
		BranchName: branchName,
	})
}

func (c *Client) DeleteBranch(branchId []byte) error {
	_, err := c.branchClient.DeleteBranch(context.Background(), &pb.DeleteBranchRequest{BranchId: branchId})
	return err
}

func (c *Client) GetBranchEncToken(branchId []byte) ([]byte, error) {
	res, err := c.branchClient.GetBranchEncToken(context.Background(), &pb.GetBranchEncTokenRequest{BranchId: branchId})
	return unmarshalStringResponse(res), err
}
