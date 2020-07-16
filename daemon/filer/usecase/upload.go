package usecase

import (
	"github.com/funnyecho/code-push/daemon/filer"
	"github.com/pkg/errors"
	"io"
)

func (c *UseCase) UploadToAliOss(stream io.Reader) (filer.FileValue, error) {
	ossKey, uploadErr := c.Adapters.Upload(stream)

	if uploadErr != nil {
		return nil, errors.WithStack(uploadErr)
	}

	return []byte(encodeAliOssObjectKey(ossKey)), nil
}
