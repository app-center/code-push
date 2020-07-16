package domain

import "github.com/funnyecho/code-push/daemon/filer"

type Service struct {
	FileService
}

type FileService interface {
	File(fileKey filer.FileKey) (*filer.File, error)
	InsertFile(file *filer.File) error
	IsFileKeyExisted(fileKey filer.FileKey) bool
}
