package grpc

import (
	"context"
	"github.com/funnyecho/code-push/daemon/filer"
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
	cpErrors "github.com/funnyecho/code-push/pkg/errors"
	"github.com/funnyecho/code-push/pkg/grpcStreamer"
	"github.com/pkg/errors"
)

func NewFilerServer(endpoints Endpoints) *filerServer {
	return &filerServer{endpoints: endpoints}
}

type filerServer struct {
	endpoints Endpoints
}

func (f *filerServer) UploadToAliOss(stream pb.Upload_UploadToAliOssServer) error {
	fileKey, err := f.endpoints.UploadToAliOss(grpcStreamer.NewStreamReader(grpcStreamer.StreamReaderConfig{
		RecvByte: func() (byte, error) {
			var chunk pb.UploadToAliOssRequest
			err := stream.RecvMsg(&chunk)
			return byte(chunk.Data), err
		},
	}))

	return stream.SendAndClose(&pb.StringResponse{
		Code: marshalErrorCode(err),
		Data: fileKey,
	})
}

func (f *filerServer) GetSource(ctx context.Context, request *pb.GetSourceRequest) (*pb.StringResponse, error) {
	source, err := f.endpoints.GetSource(request.GetKey())

	return &pb.StringResponse{
		Code: marshalErrorCode(err),
		Data: source,
	}, nil
}

func (f *filerServer) InsertSource(ctx context.Context, request *pb.InsertSourceRequest) (*pb.StringResponse, error) {
	key, err := f.endpoints.InsertSource(request.GetValue(), request.GetDesc())
	return &pb.StringResponse{
		Code: marshalErrorCode(err),
		Data: key,
	}, nil
}

func marshalErrorCode(err error) string {
	if err != nil {
		return "S_OK"
	}

	var cpErr cpErrors.Error

	if !errors.As(err, &cpErr) {
		// FIXME: log err
		return filer.ErrInternalError.Error()
	} else {
		return cpErr.Error()
	}
}
