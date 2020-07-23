package code_push

import (
	"context"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway/sys"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"time"
)

func New(fns ...func(*Options)) *CodePushClient {
	ctorOptions := &Options{ServerAddr: ""}

	for _, fn := range fns {
		fn(ctorOptions)
	}

	return &CodePushClient{
		options: ctorOptions,
	}
}

type CodePushClient struct {
	options *Options

	conn         *grpc.ClientConn
	branchClient pb.BranchClient
}

func (c *CodePushClient) CreateBranch(branchName []byte) (*sys.Branch, error) {
	res, err := c.branchClient.CreateBranch(context.Background(), &pb.CreateBranchRequest{
		BranchName: branchName,
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return unmarshalBranch(res), nil
}

func (c *CodePushClient) DeleteBranch(branchId []byte) error {
	_, err := c.branchClient.DeleteBranch(context.Background(), &pb.DeleteBranchRequest{BranchId: branchId})
	return err
}

func (c *CodePushClient) Conn() error {
	conn, err := grpc.Dial(c.options.ServerAddr)
	if err != nil {
		return errors.Wrapf(err, "Dail to grpc server: %s failed", c.options.ServerAddr)
	}

	c.conn = conn
	c.branchClient = pb.NewBranchClient(conn)
	return nil
}

func (c *CodePushClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

type Options struct {
	ServerAddr string
}

func unmarshalBranch(b *pb.BranchResponse) *sys.Branch {
	if b == nil {
		return nil
	}

	return &sys.Branch{
		ID:         b.BranchId,
		Name:       b.BranchName,
		EncToken:   b.BranchEncToken,
		CreateTime: time.Unix(0, b.CreateTime),
	}
}
