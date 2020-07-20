package internal

import (
	"github.com/funnyecho/code-push/daemon/code-push"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"time"
)

//go:generate protoc --gogofaster_out=. internal.proto

func MarshalBranch(b *code_push.Branch) (bytes []byte, err error) {
	bytes, err = proto.Marshal(&Branch{
		ID:         b.ID,
		Name:       b.Name,
		EncToken:   b.EncToken,
		CreateTime: b.CreateTime.UnixNano(),
	})

	if err != nil {
		err = errors.Wrap(err, "protobuf marshal failed")
	}

	return
}

func UnmarshalBranch(data []byte, b *code_push.Branch) error {
	var pb Branch
	if err := proto.Unmarshal(data, &pb); err != nil {
		return errors.Wrap(err, "protobuf unmarshal failed")
	}

	b.ID = pb.GetID()
	b.Name = pb.GetName()
	b.EncToken = pb.GetEncToken()
	b.CreateTime = time.Unix(0, pb.GetCreateTime()).UTC()

	return nil
}

func MarshalEnv(e *code_push.Env) (bytes []byte, err error) {
	bytes, err = proto.Marshal(&Env{
		BranchId:   e.BranchId,
		ID:         e.ID,
		Name:       e.Name,
		EncToken:   e.EncToken,
		CreateTime: e.CreateTime.UnixNano(),
	})

	if err != nil {
		err = errors.Wrap(err, "protobuf marshal failed")
	}

	return
}

func UnmarshalEnv(data []byte, e *code_push.Env) error {
	var pb Env
	if err := proto.Unmarshal(data, &pb); err != nil {
		return errors.Wrap(err, "protobuf unmarshal failed")
	}

	e.BranchId = pb.GetBranchId()
	e.ID = pb.GetID()
	e.Name = pb.GetName()
	e.EncToken = pb.GetEncToken()
	e.CreateTime = time.Unix(0, pb.GetCreateTime()).UTC()

	return nil
}

func MarshalVersion(v *code_push.Version) (bytes []byte, err error) {
	bytes, err = proto.Marshal(&Version{
		EnvId:            v.EnvId,
		AppVersion:       v.AppVersion,
		CompatAppVersion: v.CompatAppVersion,
		MustUpdate:       v.MustUpdate,
		Changelog:        v.Changelog,
		PackageFileKey:   v.PackageFileKey,
		CreateTime:       v.CreateTime.UnixNano(),
	})

	if err != nil {
		err = errors.Wrap(err, "protobuf marshal failed")
	}

	return
}

func UnmarshalVersion(data []byte, v *code_push.Version) error {
	var pb Version
	if err := proto.Unmarshal(data, &pb); err != nil {
		return errors.Wrap(err, "protobuf unmarshal failed")
	}

	v.EnvId = pb.GetEnvId()
	v.AppVersion = pb.GetAppVersion()
	v.CompatAppVersion = pb.GetCompatAppVersion()
	v.MustUpdate = pb.GetMustUpdate()
	v.Changelog = pb.GetChangelog()
	v.PackageFileKey = pb.GetPackageFileKey()

	v.CreateTime = time.Unix(0, pb.GetCreateTime()).UTC()

	return nil
}
