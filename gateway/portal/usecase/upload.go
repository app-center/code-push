package usecase

import (
	"mime/multipart"
)

func (u *useCase) UploadPkg(stream multipart.File) (fileKey []byte, err error) {
	return u.filer.UploadPkg(stream)
}
