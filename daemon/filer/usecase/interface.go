package usecase

import (
	"github.com/funnyecho/code-push/daemon/filer"
	"io"
)

type File interface {
	GetSource(key string) (*filer.File, error)
	InsertSource(value, desc, fileMD5 string, fileSize int64) (filer.FileKey, error)
}

type Upload interface {
	UploadToAliOss(stream io.Reader) (filer.FileKey, error)
}

type DomainAdapter interface {
	File(fileKey string) (*filer.File, error)
	InsertFile(file *filer.File) error
	IsFileKeyExisted(fileKey string) bool
}

type AliOssAdapter interface {
	SignFetchURL(key []byte) ([]byte, error)
	Upload(stream io.Reader) ([]byte, error)
}
