package usecase

import (
	"context"
	"mime/multipart"
)

func (u *useCase) UploadPkg(ctx context.Context, stream multipart.File) (fileKey []byte, err error) {
	return u.filer.UploadPkg(ctx, stream)
}
