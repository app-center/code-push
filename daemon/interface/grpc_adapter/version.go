package adapter

import (
	"context"
	"github.com/funnyecho/code-push/daemon/interface/grpc/pb"
)

func (c *Client) ReleaseVersion(ctx context.Context, params *pb.VersionReleaseRequest) error {
	_, err := c.versionClient.ReleaseVersion(ctx, params)
	return err
}

func (c *Client) GetVersion(ctx context.Context, envId, appVersion []byte) (*pb.VersionResponse, error) {
	return c.versionClient.GetVersion(ctx, &pb.GetVersionRequest{
		EnvId:      envId,
		AppVersion: appVersion,
	})
}

func (c *Client) GetVersionList(ctx context.Context, envId []byte) ([]*pb.VersionResponse, error) {
	res, err := c.versionClient.ListVersions(ctx, &pb.ListVersionsRequest{EnvId: envId})
	if err != nil {
		return nil, err
	}

	return res.List, nil
}

func (c *Client) VersionStrictCompatQuery(ctx context.Context, envId, appVersion []byte) (*pb.VersionStrictCompatQueryResponse, error) {
	return c.versionClient.VersionStrictCompatQuery(ctx, &pb.VersionStrictCompatQueryRequest{
		EnvId:      envId,
		AppVersion: appVersion,
	})
}
