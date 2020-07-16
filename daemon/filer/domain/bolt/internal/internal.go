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
		Key:        string(f.Key),
		Value:      string(f.Value),
		Desc:       string(f.Desc),
		CreateTime: f.CreateTime.UnixNano(),
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

	if key := pb.GetKey(); len(key) != 0 {
		f.Key = []byte(key)
	}

	if value := pb.GetValue(); len(value) != 0 {
		f.Value = []byte(value)
	}

	if desc := pb.GetDesc(); len(desc) != 0 {
		f.Desc = []byte(desc)
	}

	f.CreateTime = time.Unix(0, pb.GetCreateTime()).UTC()

	return nil
}
