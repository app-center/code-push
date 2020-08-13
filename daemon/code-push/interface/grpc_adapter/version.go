package codePushAdapter

import (
	"context"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
)

func (c *Client) ReleaseVersion(params *pb.VersionReleaseRequest) error {
	_, err := c.versionClient.ReleaseVersion(context.Background(), params)
	return err
}

func (c *Client) GetVersion(envId, appVersion []byte) (*pb.VersionResponse, error) {
	return c.versionClient.GetVersion(context.Background(), &pb.GetVersionRequest{
		EnvId:      envId,
		AppVersion: appVersion,
	})
}

func (c *Client) VersionStrictCompatQuery(envId, appVersion []byte) (*pb.VersionStrictCompatQueryResponse, error) {
	return c.versionClient.VersionStrictCompatQuery(context.Background(), &pb.VersionStrictCompatQueryRequest{
		EnvId:      envId,
		AppVersion: appVersion,
	})
}
