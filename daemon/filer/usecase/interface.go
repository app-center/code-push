package usecase

import (
	"github.com/funnyecho/code-push/daemon/filer"
	"io"
)

type Adapters struct {
	DomainAdapter
	AliOssAdapter
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
