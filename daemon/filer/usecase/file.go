package usecase

import (
	"fmt"
	"github.com/funnyecho/code-push/daemon/filer"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"net/url"
)

func (c *UseCase) InsertSource(value, desc, fileMD5 string, fileSize int64) (filer.FileKey, error) {
	if value == "" {
		return nil, errors.Wrap(filer.ErrInvalidFileValue, "filer.File value required")
	}

	u, uErr := url.Parse(string(value))
	if uErr != nil {
		return nil, errors.Wrap(filer.ErrInvalidFileValue, "filer.File value not a valid uri string")
	}

	switch u.Scheme {
	case schemeAliOss:
		fallthrough
	case schemeHttp:
		fallthrough
	case schemeHttps:
		break
	default:
		return nil, errors.Wrapf(filer.ErrInvalidFileValue, "unSupported filer.File uri scheme: %s", u.Scheme)
	}

	fileKey := generateFileKey()

	err := c.domain.InsertFile(&filer.File{
		Key:      fileKey,
		Value:    value,
		Desc:     desc,
		FileMD5:  fileMD5,
		FileSize: fileSize,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to insert filer.File")
	}

	return []byte(fileKey), nil
}

func (c *UseCase) GetSource(key string) (*filer.File, error) {
	if key == "" {
		return nil, errors.Wrap(filer.ErrInvalidFileKey, "key required")
	}

	file, fileErr := c.domain.File(key)
	if fileErr != nil {
		return nil, errors.WithStack(fileErr)
	}
	if file == nil {
		return nil, errors.Wrapf(filer.ErrFileKeyNotFound, "key: %s", key)
	}

	value := file.Value
	if value == "" {
		return nil, errors.Wrap(filer.ErrInvalidFileValue, "filer.File value missed")
	}

	u, uErr := url.Parse(value)
	if uErr != nil {
		return nil, errors.Wrap(filer.ErrInvalidFileValue, "filer.File value not a valid uri string")
	}

	var source string

	switch u.Scheme {
	case schemeAliOss:
		if aliSource, aliSourceErr := c.getAliOssSource([]byte(value)); aliSourceErr != nil {
			return nil, errors.Wrapf(aliSourceErr, "failed to get alioss source:%s", value)
		} else {
			source = string(aliSource)
		}
	case schemeHttp:
		fallthrough
	case schemeHttps:
		source = value
	default:
		return nil, errors.Wrapf(filer.ErrInvalidFileValue, "unSupported filer.File uri scheme: %s", u.Scheme)
	}

	return &filer.File{
		Key:        key,
		Value:      source,
		Desc:       file.Desc,
		CreateTime: file.CreateTime,
		FileMD5:    file.FileMD5,
		FileSize:   file.FileSize,
	}, nil
}

func (c *UseCase) getAliOssSource(fileValue []byte) ([]byte, error) {
	objectKey := decodeAliOssObjectKey(fileValue)
	return c.aliOss.SignFetchURL(objectKey)
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
