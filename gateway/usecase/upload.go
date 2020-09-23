package usecase

import (
	"context"
	"mime/multipart"
)

func (uc *useCase) UploadPkg(ctx context.Context, stream multipart.File) (fileKey []byte, err error) {
	return uc.daemon.UploadPkg(ctx, stream)
}
