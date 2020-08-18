package internal

import (
	"github.com/funnyecho/code-push/daemon/filer"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"time"
)

//go:generate protoc --gogofaster_out=. internal.proto

func MarshalFile(f *filer.File) (bytes []byte, err error) {
	bytes, err = proto.Marshal(&File{
		Key:        f.Key,
		Value:      f.Value,
		Desc:       f.Desc,
		CreateTime: f.CreateTime.UnixNano(),
		FileMD5:    f.FileMD5,
		FileSize:   int64(f.FileSize),
	})

	if err != nil {
		err = errors.Wrap(err, "protobuf marshal failed")
	}

	return
}

func UnmarshalFile(data []byte, f *filer.File) error {
	var pb File
	if err := proto.Unmarshal(data, &pb); err != nil {
		return errors.Wrap(err, "protobuf unmarshal failed")
	}

	f.Key = pb.GetKey()
	f.Value = pb.GetValue()
	f.Desc = pb.GetDesc()

	f.CreateTime = time.Unix(0, pb.GetCreateTime()).UTC()
	f.FileMD5 = pb.GetFileMD5()
	f.FileSize = pb.GetFileSize()

	return nil
}
