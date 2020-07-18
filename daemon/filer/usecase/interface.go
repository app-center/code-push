package usecase

import (
	"github.com/funnyecho/code-push/daemon/filer"
	"io"
)

type File interface {
	GetSource(key filer.FileKey) ([]byte, error)
	InsertSource(value filer.FileValue, desc filer.FileDesc) (filer.FileKey, error)
}

type Upload interface {
	UploadToAliOss(stream io.Reader) (filer.FileKey, error)
}

type DomainAdapter interface {
	File(fileKey filer.FileKey) (*filer.File, error)
	InsertFile(file *filer.File) error
	IsFileKeyExisted(fileKey filer.FileKey) bool
}

type AliOssAdapter interface {
	SignFetchURL(key []byte) ([]byte, error)
	Upload(stream io.Reader) ([]byte, error)
}
