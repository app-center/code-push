package grpc

import (
	"github.com/funnyecho/code-push/daemon/filer"
	"io"
)

type Endpoints interface {
	File
	Upload
}

type File interface {
	GetSource(key filer.FileKey) (filer.FileValue, error)
	InsertSource(value filer.FileValue, desc filer.FileDesc) (filer.FileKey, error)
}

type Upload interface {
	UploadToAliOss(stream io.Reader) (filer.FileValue, error)
}
