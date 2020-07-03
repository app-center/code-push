package usecase

import (
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc"
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
	"github.com/funnyecho/code-push/daemon/filer/usecase/internal"
	"github.com/funnyecho/code-push/pkg/grpcStreamer"
)

type IUpload interface {
	UploadToAliOss(server pb.Upload_UploadToAliOssServer) (FileValue, error)
}

func NewUploadUseCase(config UploadUseCaseConfig) IUpload {
	return &uploadUseCase{aliOssClient: config.AliOssClient}
}

type uploadUseCase struct {
	aliOssClient internal.AliOssClient
}

func (u *uploadUseCase) UploadToAliOss(server pb.Upload_UploadToAliOssServer) (FileValue, error) {
	bucket, bucketErr := u.aliOssClient.GetPackageBucket()
	if bucketErr != nil {
		server.SendAndClose(&pb.StringResponse{
			Code: grpc.MarshalErrorCode(bucketErr),
			Data: "",
		})
	}

	objectKey := u.aliOssClient.GeneratePackageObjectKey()

	uploadErr := bucket.PutObject(
		objectKey,
		grpcStreamer.NewStreamReader(grpcStreamer.StreamReaderConfig{
			RecvByte: func() (byte, error) {
				var chunk pb.UploadToAliOssRequest
				err := server.RecvMsg(&chunk)
				return byte(chunk.Data), err
			},
		}),
	)

	if uploadErr != nil {
		return nil, uploadErr
	}

	return []byte(encodeAliOssObjectKey([]byte(objectKey))), nil
}

type UploadUseCaseConfig struct {
	AliOssClient internal.AliOssClient
}

type IUploadChunk struct {
	Data uint32
}
