package usecase

import (
	"fmt"
	"github.com/funnyecho/code-push/daemon"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"net/url"
)

func (uc *useCase) InsertSource(value, desc, fileMD5 string, fileSize int64) (daemon.FileKey, error) {
	if value == "" {
		return nil, errors.Wrap(daemon.ErrInvalidFileValue, "daemon.File value required")
	}

	u, uErr := url.Parse(string(value))
	if uErr != nil {
		return nil, errors.Wrap(daemon.ErrInvalidFileValue, "daemon.File value not a valid uri string")
	}

	switch u.Scheme {
	case schemeAliOss:
		fallthrough
	case schemeHttp:
		fallthrough
	case schemeHttps:
		break
	default:
		return nil, errors.Wrapf(daemon.ErrInvalidFileValue, "unSupported daemon.File uri scheme: %s", u.Scheme)
	}

	fileKey := generateFileKey()

	err := uc.domain.InsertFile(&daemon.File{
		Key:      fileKey,
		Value:    value,
		Desc:     desc,
		FileMD5:  fileMD5,
		FileSize: fileSize,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to insert daemon.File")
	}

	return []byte(fileKey), nil
}

func (uc *useCase) GetSource(key string) (*daemon.File, error) {
	if key == "" {
		return nil, errors.Wrap(daemon.ErrInvalidFileKey, "key required")
	}

	file, fileErr := uc.domain.File(key)
	if fileErr != nil {
		return nil, errors.WithStack(fileErr)
	}
	if file == nil {
		return nil, errors.Wrapf(daemon.ErrFileKeyNotFound, "key: %s", key)
	}

	value := file.Value
	if value == "" {
		return nil, errors.Wrap(daemon.ErrInvalidFileValue, "daemon.File value missed")
	}

	u, uErr := url.Parse(value)
	if uErr != nil {
		return nil, errors.Wrap(daemon.ErrInvalidFileValue, "daemon.File value not a valid uri string")
	}

	var source string

	switch u.Scheme {
	case schemeAliOss:
		if aliSource, aliSourceErr := uc.getAliOssSource([]byte(value)); aliSourceErr != nil {
			return nil, errors.Wrapf(aliSourceErr, "failed to get alioss source:%s", value)
		} else {
			source = string(aliSource)
		}
	case schemeHttp:
		fallthrough
	case schemeHttps:
		source = value
	default:
		return nil, errors.Wrapf(daemon.ErrInvalidFileValue, "unSupported daemon.File uri scheme: %s", u.Scheme)
	}

	return &daemon.File{
		Key:        key,
		Value:      source,
		Desc:       file.Desc,
		CreateTime: file.CreateTime,
		FileMD5:    file.FileMD5,
		FileSize:   file.FileSize,
	}, nil
}

func (uc *useCase) getAliOssSource(fileValue []byte) ([]byte, error) {
	objectKey := decodeAliOssObjectKey(fileValue)
	return uc.aliOss.SignFetchURL(objectKey)
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
