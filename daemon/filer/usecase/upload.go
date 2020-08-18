package usecase

import (
	"bytes"
	"github.com/funnyecho/code-push/daemon/filer"
	"github.com/funnyecho/code-push/pkg/util"
	"github.com/pkg/errors"
	"io"
)

func (c *UseCase) UploadToAliOss(stream io.Reader) (filer.FileKey, error) {
	if stream == nil {
		return nil, errors.Wrap(filer.ErrParamsInvalid, "upload stream required")
	}

	var buf bytes.Buffer
	tee := io.TeeReader(stream, &buf)

	ossKey, uploadErr := c.aliOss.Upload(tee)

	if uploadErr != nil {
		return nil, errors.WithStack(uploadErr)
	}

	fileValue := encodeAliOssObjectKey(ossKey)
	fileSize := buf.Len()
	fileMD5 := util.EncodeMD5(string(buf.Bytes()))

	return c.InsertSource(fileValue, "", fileMD5, int64(fileSize))
}
