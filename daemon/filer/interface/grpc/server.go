package grpc

import (
	"context"
	"github.com/funnyecho/code-push/daemon/filer"
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
	"github.com/funnyecho/code-push/pkg/grpc-streamer"
	"github.com/funnyecho/code-push/pkg/log"
	"io"
)

func NewFilerServer(endpoints Endpoints, logger log.Logger) *filerServer {
	return &filerServer{
		endpoints: endpoints,
		Logger:    logger,
	}
}

type filerServer struct {
	endpoints Endpoints
	log.Logger
}

func (f *filerServer) UploadToAliOss(stream pb.Upload_UploadToAliOssServer) error {
	fileKey, err := f.endpoints.UploadToAliOss(grpc_streamer.NewStreamReader(grpc_streamer.StreamReaderConfig{
		RecvByte: func() (b byte, err error) {
			chunk, recvErr := stream.Recv()
			if recvErr != nil {
				err = recvErr
				return
			}

			if chunk == nil {
				err = io.EOF
				return
			}

			return byte(chunk.Data), nil
		},
	}))

	if err != nil {
		return err
	}

	return stream.SendAndClose(marshalBytesToStringResponse(fileKey))
}

func (f *filerServer) GetSource(ctx context.Context, request *pb.GetSourceRequest) (*pb.FileSource, error) {
	source, err := f.endpoints.GetSource(request.GetKey())

	return marshalFileSource(source), err
}

func (f *filerServer) InsertSource(ctx context.Context, request *pb.InsertSourceRequest) (*pb.StringResponse, error) {
	key, err := f.endpoints.InsertSource(request.GetValue(), request.GetDesc(), request.GetFileMD5(), request.GetFileSize())
	return marshalBytesToStringResponse(key), err
}

func marshalBytesToStringResponse(data []byte) *pb.StringResponse {
	if data == nil {
		return nil
	}

	return &pb.StringResponse{Data: string(data)}
}

func marshalFileSource(file *filer.File) *pb.FileSource {
	if file == nil {
		return nil
	}

	return &pb.FileSource{
		Key:        file.Key,
		Value:      file.Value,
		Desc:       file.Desc,
		CreateTime: file.CreateTime.UnixNano(),
		FileMD5:    file.FileMD5,
		FileSize:   file.FileSize,
	}
}
