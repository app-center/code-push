package grpc

import (
	"context"
	"github.com/funnyecho/code-push/daemon/filer"
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
	cpErrors "github.com/funnyecho/code-push/pkg/errors"
	"github.com/funnyecho/code-push/pkg/grpcStreamer"
	"github.com/pkg/errors"
)

type FilerServer struct {
	Endpoints Endpoints
}

func (f *FilerServer) UploadToAliOss(stream pb.Upload_UploadToAliOssServer) error {
	fileKey, err := f.Endpoints.UploadToAliOss(grpcStreamer.NewStreamReader(grpcStreamer.StreamReaderConfig{
		RecvByte: func() (byte, error) {
			var chunk pb.UploadToAliOssRequest
			err := stream.RecvMsg(&chunk)
			return byte(chunk.Data), err
		},
	}))

	return stream.SendAndClose(&pb.StringResponse{
		Code: marshalErrorCode(err),
		Data: string(fileKey),
	})
}

func (f *FilerServer) GetSource(ctx context.Context, request *pb.GetSourceRequest) (*pb.StringResponse, error) {
	source, err := f.Endpoints.GetSource([]byte(request.GetKey()))

	return &pb.StringResponse{
		Code: marshalErrorCode(err),
		Data: string(source),
	}, nil
}

func (f *FilerServer) InsertSource(ctx context.Context, request *pb.InsertSourceRequest) (*pb.StringResponse, error) {
	key, err := f.Endpoints.InsertSource([]byte(request.GetValue()), []byte(request.GetDesc()))
	return &pb.StringResponse{
		Code: marshalErrorCode(err),
		Data: string(key),
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
