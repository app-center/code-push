package code_push

import (
	"context"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway/sys"
	"github.com/funnyecho/code-push/pkg/grpcInterceptor"
	"github.com/funnyecho/code-push/pkg/log"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"time"
)

func New(logger log.Logger, fns ...func(*Options)) *CodePushClient {
	ctorOptions := &Options{ServerAddr: ""}

	for _, fn := range fns {
		fn(ctorOptions)
	}

	return &CodePushClient{
		Logger:  logger,
		Options: ctorOptions,
	}
}

type CodePushClient struct {
	log.Logger
	*Options

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
	conn, err := grpc.Dial(
		c.Options.ServerAddr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpcInterceptor.UnaryClientMetricInterceptor(c.Logger),
			grpcInterceptor.UnaryClientErrorInterceptor(),
		)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpcInterceptor.StreamClientMetricInterceptor(c.Logger),
			grpcInterceptor.StreamClientErrorInterceptor(),
		)),
	)
	if err != nil {
		return errors.Wrapf(err, "Dail to grpc server: %s failed", c.ServerAddr)
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
