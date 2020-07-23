package grpc

import (
	"context"
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
	"github.com/funnyecho/code-push/pkg/grpcStreamer"
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

	if err != nil {
		return err
	}

	return stream.SendAndClose(marshalBytesToStringResponse(fileKey))
}

func (f *filerServer) GetSource(ctx context.Context, request *pb.GetSourceRequest) (*pb.StringResponse, error) {
	source, err := f.endpoints.GetSource(request.GetKey())

	return marshalBytesToStringResponse(source), err
}

func (f *filerServer) InsertSource(ctx context.Context, request *pb.InsertSourceRequest) (*pb.StringResponse, error) {
	key, err := f.endpoints.InsertSource(request.GetValue(), request.GetDesc())
	return marshalBytesToStringResponse(key), err
}

func marshalBytesToStringResponse(data []byte) *pb.StringResponse {
	if data == nil {
		return nil
	}

	return &pb.StringResponse{Data: string(data)}
}
