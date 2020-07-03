package domain

import "time"

type FileKey []byte

type File struct {
	Key        []byte
	Value      []byte
	Desc       []byte
	CreateTime time.Time
}

type AliOssScheme struct {
	Endpoint        []byte
	AccessKeyId     []byte
	AccessKeySecret []byte
	UpdateTime      time.Time
}

type IFileService interface {
	File(fileKey FileKey) (*File, error)
	InsertFile(file *File) error
	IsFileKeyExisted(fileKey FileKey) bool
}

type ISchemeService interface {
	AliOssScheme() (*AliOssScheme, error)
	UpdateAliOssScheme(scheme *AliOssScheme) error
}
