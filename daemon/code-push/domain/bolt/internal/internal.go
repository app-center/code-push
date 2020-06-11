package internal

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain"
	"github.com/gogo/protobuf/proto"
	"time"
)

//go:generate protoc --gogofaster_out=. internal.proto

func MarshalBranch(b *domain.Branch) ([]byte, error) {
	return proto.Marshal(&Branch{
		ID:         b.ID,
		Name:       b.Name,
		AuthHost:   b.AuthHost,
		EncToken:   b.EncToken,
		CreateTime: b.CreateTime.UnixNano(),
	})
}

func UnmarshalBranch(data []byte, b *domain.Branch) error {
	var pb Branch
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}

	b.ID = pb.GetID()
	b.Name = pb.GetName()
	b.AuthHost = pb.GetAuthHost()
	b.EncToken = pb.GetEncToken()
	b.CreateTime = time.Unix(0, pb.GetCreateTime()).UTC()

	return nil
}

func MarshalEnv(e *domain.Env) ([]byte, error) {
	return proto.Marshal(&Env{
		BranchId:   e.BranchId,
		ID:         e.ID,
		Name:       e.Name,
		EncToken:   e.EncToken,
		CreateTime: e.CreateTime.UnixNano(),
	})
}

func UnmarshalEnv(data []byte, e *domain.Env) error {
	var pb Env
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}

	e.BranchId = pb.GetBranchId()
	e.ID = pb.GetID()
	e.Name = pb.GetName()
	e.EncToken = pb.GetEncToken()
	e.CreateTime = time.Unix(0, pb.GetCreateTime()).UTC()

	return nil
}

func MarshalVersion(v *domain.Version) ([]byte, error) {
	return proto.Marshal(&Version{
		EnvId:            v.EnvId,
		AppVersion:       v.AppVersion,
		CompatAppVersion: v.CompatAppVersion,
		MustUpdate:       v.MustUpdate,
		Changelog:        v.Changelog,
		PackageUri:       v.PackageUri,
		CreateTime:       v.CreateTime.UnixNano(),
	})
}

func UnmarshalVersion(data []byte, v *domain.Version) error {
	var pb Version
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}

	v.EnvId = pb.GetEnvId()
	v.AppVersion = pb.GetAppVersion()
	v.CompatAppVersion = pb.GetCompatAppVersion()
	v.MustUpdate = pb.GetMustUpdate()
	v.Changelog = pb.GetChangelog()
	v.PackageUri = pb.GetPackageUri()

	v.CreateTime = time.Unix(0, pb.GetCreateTime()).UTC()

	return nil
}