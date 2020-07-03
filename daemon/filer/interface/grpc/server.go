package grpc

import (
	"context"
	code_push "github.com/funnyecho/code-push/daemon/code-push"
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
	"github.com/funnyecho/code-push/daemon/filer/usecase"
	cpErrors "github.com/funnyecho/code-push/pkg/errors"
	"github.com/pkg/errors"
)

func NewFilerServer(config FilerServerConfig) *filerServer {
	return &filerServer{
		fileUseCase:   config.FileUseCase,
		schemeUseCase: config.SchemeUseCase,
		uploadUseCase: config.UploadUseCase,
	}
}

type filerServer struct {
	fileUseCase   usecase.IFile
	schemeUseCase usecase.IScheme
	uploadUseCase usecase.IUpload
}

func (f *filerServer) UploadToAliOss(stream pb.Upload_UploadToAliOssServer) error {
	fileKey, err := f.uploadUseCase.UploadToAliOss(stream)
	stream.SendAndClose(&pb.StringResponse{
		Code: MarshalErrorCode(err),
		Data: string(fileKey),
	})

	return nil
}

func (f *filerServer) UpdateAliOssScheme(ctx context.Context, config *pb.AliOssSchemeConfig) (*pb.PlainResponse, error) {
	updateConfig := &usecase.AliOssSchemeConfig{
		Endpoint:        nil,
		AccessKeyId:     nil,
		AccessKeySecret: nil,
	}

	if len(config.Endpoint) > 0 {
		updateConfig.Endpoint = []byte(config.Endpoint)
	}

	if len(config.AccessKeyId) > 0 {
		updateConfig.AccessKeyId = []byte(config.AccessKeyId)
	}

	if len(config.AccessKeySecret) > 0 {
		updateConfig.AccessKeySecret = []byte(config.AccessKeySecret)
	}

	err := f.schemeUseCase.UpdateAliOssScheme(updateConfig)
	return &pb.PlainResponse{Code: MarshalErrorCode(err)}, nil
}

func (f *filerServer) GetSource(ctx context.Context, request *pb.GetSourceRequest) (*pb.StringResponse, error) {
	source, err := f.fileUseCase.GetSource([]byte(request.GetKey()))

	return &pb.StringResponse{
		Code: MarshalErrorCode(err),
		Data: string(source),
	}, nil
}

func (f *filerServer) InsertSource(ctx context.Context, request *pb.InsertSourceRequest) (*pb.StringResponse, error) {
	key, err := f.fileUseCase.InsertSource([]byte(request.GetValue()), []byte(request.GetDesc()))
	return &pb.StringResponse{
		Code: MarshalErrorCode(err),
		Data: string(key),
	}, nil
}

func MarshalErrorCode(err error) string {
	if err != nil {
		return "S_OK"
	}

	var cpErr cpErrors.Error

	if !errors.As(err, &cpErr) {
		// FIXME: log err
		return code_push.ErrInternalError.Error()
	} else {
		return cpErr.Error()
	}
}

type FilerServerConfig struct {
	FileUseCase   usecase.IFile
	SchemeUseCase usecase.IScheme
	UploadUseCase usecase.IUpload
}
