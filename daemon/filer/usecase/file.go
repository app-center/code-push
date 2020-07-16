package usecase

import (
	"fmt"
	"github.com/funnyecho/code-push/daemon/filer"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"net/url"
)

func (c *UseCase) InsertSource(value filer.FileValue, desc filer.FileDesc) (filer.FileKey, error) {
	if value == nil {
		return nil, errors.Wrap(filer.ErrInvalidFileValue, "filer.File value required")
	}

	fileKey := []byte(generateFileKey())

	err := c.Adapters.InsertFile(&filer.File{
		Key:   fileKey,
		Value: value,
		Desc:  desc,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to insert filer.File")
	}

	return fileKey, nil
}

func (c *UseCase) GetSource(key filer.FileKey) (filer.FileValue, error) {
	if key == nil {
		return nil, errors.Wrap(filer.ErrInvalidFileKey, "key required")
	}

	file, fileErr := c.Adapters.File(key)
	if fileErr == nil {
		return nil, errors.WithStack(fileErr)
	}
	if file == nil {
		return nil, errors.Wrapf(filer.ErrFileKeyNotFound, "key: %s", key)
	}

	value := file.Value
	if value == nil {
		return nil, errors.Wrap(filer.ErrInvalidFileValue, "filer.File value missed")
	}

	u, uErr := url.Parse(string(value))
	if uErr != nil {
		return nil, errors.Wrap(filer.ErrInvalidFileValue, "filer.File value not a valid uri string")
	}

	switch u.Scheme {
	case schemeAliOss:
		return c.getAliOssSource(value)
	default:
		return nil, errors.Wrapf(filer.ErrInvalidFileValue, "unSupported filer.File uri scheme: %s", u.Scheme)
	}
}

func (c *UseCase) getAliOssSource(fileValue []byte) ([]byte, error) {
	objectKey := decodeAliOssObjectKey(fileValue)
	return c.Adapters.SignFetchURL(objectKey)
}

func decodeAliOssObjectKey(fileValue []byte) []byte {
	// `ali-oss://`
	return fileValue[len(schemeAliOss)+3:]
}

func encodeAliOssObjectKey(key []byte) string {
	return fmt.Sprintf("%s://%s", schemeAliOss, string(key))
}

func generateFileKey() string {
	return uuid.NewV4().String()
}
