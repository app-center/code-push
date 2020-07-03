package internal

import (
	"github.com/funnyecho/code-push/daemon/filer/domain"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"time"
)

//go:generate protoc --gogofaster_out=. internal.proto

func MarshalFile(f *domain.File) (bytes []byte, err error) {
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

func UnmarshalFile(data []byte, f *domain.File) error {
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

func MarshalAliOssScheme(s *domain.AliOssScheme) (bytes []byte, err error) {
	bytes, err = proto.Marshal(&AliOssScheme{
		Endpoint:        string(s.Endpoint),
		AccessKeyId:     string(s.AccessKeyId),
		AccessKeySecret: string(s.AccessKeySecret),
		UpdateTime:      s.UpdateTime.UnixNano(),
	})

	if err != nil {
		err = errors.Wrap(err, "protobuf marshal failed")
	}

	return
}

func UnmarshalAliOssScheme(data []byte, s *domain.AliOssScheme) error {
	var pb AliOssScheme
	if err := proto.Unmarshal(data, &pb); err != nil {
		return errors.Wrap(err, "protobuf unmarshal failed")
	}

	s.Endpoint = []byte(pb.GetEndpoint())
	s.AccessKeyId = []byte(pb.GetAccessKeyId())
	s.AccessKeySecret = []byte(pb.GetAccessKeySecret())
	s.UpdateTime = time.Unix(0, pb.GetUpdateTime()).UTC()

	return nil
}
