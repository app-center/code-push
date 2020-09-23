package usecase

import (
	"bytes"
	"github.com/funnyecho/code-push/daemon"
	"github.com/funnyecho/code-push/pkg/util"
	"github.com/pkg/errors"
	"io"
)

func (uc *useCase) UploadToAliOss(stream io.Reader) (daemon.FileKey, error) {
	if stream == nil {
		return nil, errors.Wrap(daemon.ErrParamsInvalid, "upload stream required")
	}

	var buf bytes.Buffer
	tee := io.TeeReader(stream, &buf)

	ossKey, uploadErr := uc.aliOss.Upload(tee)

	if uploadErr != nil {
		return nil, errors.WithStack(uploadErr)
	}

	fileValue := encodeAliOssObjectKey(ossKey)
	fileSize := buf.Len()
	fileMD5 := util.EncodeMD5(string(buf.Bytes()))

	return uc.InsertSource(fileValue, "", fileMD5, int64(fileSize))
}
