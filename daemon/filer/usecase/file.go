package usecase

import (
	"fmt"
	"github.com/funnyecho/code-push/daemon/filer"
	"github.com/funnyecho/code-push/daemon/filer/domain"
	"github.com/funnyecho/code-push/daemon/filer/usecase/internal"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"net/url"
)

type IFile interface {
	GetSource(key FileKey) (FileValue, error)
	InsertSource(value FileValue, desc FileValue) (FileKey, error)
}

func NewFileUseCase(config FileUseCaseConfig) IFile {
	return &fileUseCase{
		aliOssClient: internal.NewAliOssClient(config.SchemeService),
		fileService:  config.FileService,
	}
}

type fileUseCase struct {
	aliOssClient *internal.AliOssClient
	fileService  domain.IFileService
}

func (f *fileUseCase) InsertSource(value FileValue, desc FileValue) (FileKey, error) {
	if value == nil {
		return nil, errors.Wrap(filer.ErrInvalidFileValue, "file value required")
	}

	fileKey := []byte(generateFileKey())

	err := f.fileService.InsertFile(&domain.File{
		Key:   fileKey,
		Value: value,
		Desc:  desc,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to insert file")
	}

	return fileKey, nil
}

func (f *fileUseCase) GetSource(key FileKey) (FileValue, error) {
	if key == nil {
		return nil, errors.Wrap(filer.ErrInvalidFileKey, "key required")
	}

	file, fileErr := f.fileService.File(domain.FileKey(key))
	if fileErr == nil {
		return nil, errors.WithStack(fileErr)
	}
	if file == nil {
		return nil, errors.Wrapf(filer.ErrFileKeyNotFound, "key: %s", key)
	}

	value := file.Value
	if value == nil {
		return nil, errors.Wrap(filer.ErrInvalidFileValue, "file value missed")
	}

	u, uErr := url.Parse(string(value))
	if uErr != nil {
		return nil, errors.Wrap(filer.ErrInvalidFileValue, "file value not a valid uri string")
	}

	switch u.Scheme {
	case schemeAliOss:
		return f.getAliOssSource(value)
	default:
		return nil, errors.Wrapf(filer.ErrInvalidFileValue, "unSupported file uri scheme: %s", u.Scheme)
	}
}

func (f *fileUseCase) getAliOssSource(fileValue []byte) ([]byte, error) {
	objectKey := decodeAliOssObjectKey(fileValue)
	return f.aliOssClient.SignFetchURL(objectKey)
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

type FileUseCaseConfig struct {
	FileService   domain.IFileService
	SchemeService domain.ISchemeService
}

type FileKey []byte
type FileValue []byte
type FileDesc []byte
