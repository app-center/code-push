package usecase

import (
	"context"
	"github.com/funnyecho/code-push/daemon/interface/grpc/pb"
	"mime/multipart"
)

type DaemonAdapter interface {
	CreateBranch(ctx context.Context, branchName []byte) (*pb.BranchResponse, error)
	DeleteBranch(ctx context.Context, branchId []byte) error
	GetBranchEncToken(ctx context.Context, branchId []byte) ([]byte, error)

	CreateEnv(ctx context.Context, branchId, envName []byte) (*pb.EnvResponse, error)
	GetEnv(ctx context.Context, envId []byte) (*pb.EnvResponse, error)
	DeleteEnv(ctx context.Context, envId []byte) error
	GetEnvEncToken(ctx context.Context, envId []byte) ([]byte, error)

	ReleaseVersion(ctx context.Context, params *pb.VersionReleaseRequest) error
	GetVersion(ctx context.Context, envId, appVersion []byte) (*pb.VersionResponse, error)
	VersionStrictCompatQuery(ctx context.Context, envId, appVersion []byte) (*pb.VersionStrictCompatQueryResponse, error)

	GenerateAccessToken(ctx context.Context, issuer pb.AccessTokenIssuer, subject string) ([]byte, error)
	VerifyAccessToken(ctx context.Context, token string) (subject []byte, err error)

	UploadPkg(ctx context.Context, source multipart.File) (fileKey []byte, err error)

	GetSource(ctx context.Context, fileKey []byte) (*pb.FileSource, error)
}