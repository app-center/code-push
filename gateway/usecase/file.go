package usecase

import (
	"context"
	"github.com/funnyecho/code-push/daemon/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway"
	"time"
)

func (uc *useCase) FileDownload(ctx context.Context, fileId []byte) ([]byte, error) {
	source, sourceErr := uc.daemon.GetSource(ctx, fileId)
	if sourceErr != nil {
		return nil, sourceErr
	}

	return []byte(source.GetValue()), nil
}

func unmarshalFileSource(v *pb.FileSource) *gateway.FileSource {
	if v == nil {
		return nil
	}

	return &gateway.FileSource{
		Key:        v.GetKey(),
		Value:      v.GetValue(),
		Desc:       v.GetDesc(),
		CreateTime: time.Unix(0, v.GetCreateTime()),
		FileMD5:    v.GetFileMD5(),
		FileSize:   v.GetFileSize(),
	}
}
