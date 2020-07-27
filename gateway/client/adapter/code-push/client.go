package code_push

import (
	"context"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway/client"
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

	conn          *grpc.ClientConn
	envClient     pb.EnvClient
	versionClient pb.VersionClient
}

func (c *CodePushClient) GetEnvEncToken(envId []byte) ([]byte, error) {
	res, err := c.envClient.GetEnvEncToken(context.Background(), &pb.EnvIdRequest{EnvId: envId})
	return unmarshalStringResponse(res), err
}

func (c *CodePushClient) GetVersion(envId, appVersion []byte) (*client.Version, error) {
	res, err := c.versionClient.GetVersion(context.Background(), &pb.GetVersionRequest{
		EnvId:      envId,
		AppVersion: appVersion,
	})

	return unmarshalVersion(res), err
}

func (c *CodePushClient) VersionStrictCompatQuery(envId, appVersion []byte) (*client.VersionCompatQueryResult, error) {
	res, err := c.versionClient.VersionStrictCompatQuery(context.Background(), &pb.VersionStrictCompatQueryRequest{
		EnvId:      envId,
		AppVersion: appVersion,
	})

	return unmarshalVersionCompatQueryResult(res), err
}

func (c *CodePushClient) Conn() error {
	conn, err := grpc.Dial(c.options.ServerAddr)
	if err != nil {
		return errors.Wrapf(err, "Dail to grpc server: %s failed", c.options.ServerAddr)
	}

	c.conn = conn
	c.envClient = pb.NewEnvClient(conn)
	c.versionClient = pb.NewVersionClient(conn)
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

func unmarshalVersion(v *pb.VersionResponse) *client.Version {
	if v == nil {
		return nil
	}

	return &client.Version{
		EnvId:            v.GetEnvId(),
		AppVersion:       v.GetAppVersion(),
		CompatAppVersion: v.GetCompatAppVersion(),
		MustUpdate:       v.GetMustUpdate(),
		Changelog:        v.GetChangelog(),
		PackageFileKey:   v.GetPackageFileKey(),
		CreateTime:       time.Unix(0, v.GetCreateTime()),
	}
}

func unmarshalVersionCompatQueryResult(r *pb.VersionStrictCompatQueryResponse) *client.VersionCompatQueryResult {
	if r == nil {
		return nil
	}

	return &client.VersionCompatQueryResult{
		AppVersion:          r.GetAppVersion(),
		LatestAppVersion:    r.GetLatestAppVersion(),
		CanUpdateAppVersion: r.GetCanUpdateAppVersion(),
		MustUpdate:          r.GetMustUpdate(),
	}
}

func unmarshalStringResponse(r *pb.StringResponse) []byte {
	if r == nil {
		return nil
	}

	return []byte(r.Data)
}
