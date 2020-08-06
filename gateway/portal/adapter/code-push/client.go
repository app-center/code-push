package code_push

import (
	"context"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway/portal"
	"github.com/funnyecho/code-push/pkg/log"
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

	conn          *grpc.ClientConn
	branchClient  pb.BranchClient
	envClient     pb.EnvClient
	versionClient pb.VersionClient
}

func (c *CodePushClient) GetBranchEncToken(branchId []byte) ([]byte, error) {
	res, err := c.branchClient.GetBranchEncToken(context.Background(), &pb.GetBranchEncTokenRequest{BranchId: branchId})
	return unmarshalStringResponse(res), err
}

func (c *CodePushClient) CreateEnv(branchId, envName []byte) (*portal.Env, error) {
	res, err := c.envClient.CreateEnv(context.Background(), &pb.CreateEnvRequest{
		BranchId: branchId,
		EnvName:  envName,
	})

	return unmarshalEnv(res), err
}

func (c *CodePushClient) GetEnv(envId []byte) (*portal.Env, error) {
	res, err := c.envClient.GetEnv(context.Background(), &pb.EnvIdRequest{EnvId: envId})
	return unmarshalEnv(res), err
}

func (c *CodePushClient) DeleteEnv(envId []byte) error {
	_, err := c.envClient.DeleteEnv(context.Background(), &pb.EnvIdRequest{EnvId: envId})
	return err
}

func (c *CodePushClient) GetEnvEncToken(envId []byte) ([]byte, error) {
	res, err := c.envClient.GetEnvEncToken(context.Background(), &pb.EnvIdRequest{EnvId: envId})
	return unmarshalStringResponse(res), err
}

func (c *CodePushClient) ReleaseVersion(params *portal.VersionReleaseParams) error {
	_, err := c.versionClient.ReleaseVersion(context.Background(), marshalVersionReleaseParams(params))
	return err
}

func (c *CodePushClient) GetVersion(envId, appVersion []byte) (*portal.Version, error) {
	res, err := c.versionClient.GetVersion(context.Background(), &pb.GetVersionRequest{
		EnvId:      envId,
		AppVersion: appVersion,
	})

	return unmarshalVersion(res), err
}

func (c *CodePushClient) VersionStrictCompatQuery(envId, appVersion []byte) (*portal.VersionCompatQueryResult, error) {
	res, err := c.versionClient.VersionStrictCompatQuery(context.Background(), &pb.VersionStrictCompatQueryRequest{
		EnvId:      envId,
		AppVersion: appVersion,
	})

	return unmarshalVersionCompatQueryResult(res), err
}

func (c *CodePushClient) Conn() error {
	conn, err := grpc.Dial(c.Options.ServerAddr, grpc.WithInsecure())
	if err != nil {
		return errors.Wrapf(err, "Dail to grpc server: %s failed", c.Options.ServerAddr)
	}

	c.conn = conn
	c.branchClient = pb.NewBranchClient(conn)
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

func unmarshalEnv(e *pb.EnvResponse) *portal.Env {
	if e == nil {
		return nil
	}

	return &portal.Env{
		BranchId:   e.GetBranchId(),
		ID:         e.GetEnvId(),
		Name:       e.GetName(),
		EncToken:   e.GetEnvEncToken(),
		CreateTime: time.Unix(0, e.CreateTime),
	}
}

func unmarshalVersion(v *pb.VersionResponse) *portal.Version {
	if v == nil {
		return nil
	}

	return &portal.Version{
		EnvId:            v.GetEnvId(),
		AppVersion:       v.GetAppVersion(),
		CompatAppVersion: v.GetCompatAppVersion(),
		MustUpdate:       v.GetMustUpdate(),
		Changelog:        v.GetChangelog(),
		PackageFileKey:   v.GetPackageFileKey(),
		CreateTime:       time.Unix(0, v.GetCreateTime()),
	}
}

func unmarshalVersionCompatQueryResult(r *pb.VersionStrictCompatQueryResponse) *portal.VersionCompatQueryResult {
	if r == nil {
		return nil
	}

	return &portal.VersionCompatQueryResult{
		AppVersion:          r.GetAppVersion(),
		LatestAppVersion:    r.GetLatestAppVersion(),
		CanUpdateAppVersion: r.GetCanUpdateAppVersion(),
		MustUpdate:          r.GetMustUpdate(),
	}
}

func marshalVersionReleaseParams(p *portal.VersionReleaseParams) *pb.VersionReleaseRequest {
	if p == nil {
		return nil
	}

	return &pb.VersionReleaseRequest{
		EnvId:            p.EnvId,
		AppVersion:       p.AppVersion,
		CompatAppVersion: p.CompatAppVersion,
		Changelog:        p.Changelog,
		PackageFileKey:   p.PackageFileKey,
		MustUpdate:       p.MustUpdate,
	}
}

func unmarshalStringResponse(r *pb.StringResponse) []byte {
	if r == nil {
		return nil
	}

	return []byte(r.Data)
}
