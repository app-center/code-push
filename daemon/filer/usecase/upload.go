package usecase

import (
	"github.com/funnyecho/code-push/daemon/filer"
	"github.com/pkg/errors"
	"io"
)

func (c *UseCase) UploadToAliOss(stream io.Reader) (filer.FileKey, error) {
	if stream == nil {
		return nil, errors.Wrap(filer.ErrParamsInvalid, "upload stream required")
	}

	ossKey, uploadErr := c.aliOss.Upload(stream)

	if uploadErr != nil {
		return nil, errors.WithStack(uploadErr)
	}

	fileValue := []byte(encodeAliOssObjectKey(ossKey))

	return c.InsertSource(fileValue, nil)
}
