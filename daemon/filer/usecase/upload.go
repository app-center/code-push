package usecase

import (
	"github.com/funnyecho/code-push/daemon/filer/domain"
	"github.com/funnyecho/code-push/daemon/filer/usecase/internal"
	"io"
)

type IUpload interface {
	UploadToAliOss(stream io.Reader) (FileValue, error)
}

func NewUploadUseCase(config UploadUseCaseConfig) IUpload {
	return &uploadUseCase{
		aliOssClient: internal.NewAliOssClient(config.SchemeService),
	}
}

type uploadUseCase struct {
	aliOssClient *internal.AliOssClient
}

func (u *uploadUseCase) UploadToAliOss(stream io.Reader) (FileValue, error) {
	bucket, bucketErr := u.aliOssClient.GetPackageBucket()
	if bucketErr != nil {
		return nil, bucketErr
	}

	objectKey := u.aliOssClient.GeneratePackageObjectKey()

	uploadErr := bucket.PutObject(
		objectKey,
		stream,
	)

	if uploadErr != nil {
		return nil, uploadErr
	}

	return []byte(encodeAliOssObjectKey([]byte(objectKey))), nil
}

type UploadUseCaseConfig struct {
	SchemeService domain.ISchemeService
}

type IUploadChunk struct {
	Data uint32
}
