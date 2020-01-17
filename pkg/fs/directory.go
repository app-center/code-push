package fs

import (
	"github.com/funnyecho/code-push/pkg/fs/errors"
	"os"
	"path/filepath"
)

type Directory struct {
	dirPath string
}

func (d *Directory) Path() string {
	return d.dirPath
}

func (d *Directory) CheckNotExist() bool {
	_, err := os.Stat(d.dirPath)
	return os.IsNotExist(err)
}

func (d *Directory) CheckPermissionDenied() bool {
	_, err := os.Stat(d.dirPath)
	return os.IsPermission(err)
}

func (d *Directory) EnsurePath() error {
	if d.CheckNotExist() {
		if err := os.MkdirAll(d.dirPath, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (d *Directory) Delete() error {
	if d.CheckNotExist() {
		return nil
	}

	return os.RemoveAll(d.Path())
}

type DirectoryConfig struct {
	DirPath string
}

func NewDirectory(config DirectoryConfig) (dir *Directory, err error) {
	if len(config.DirPath) <= 0 {
		return nil, errors.NewInvalidPathError(errors.InvalidPathConfig{
			Path: config.DirPath,
		})
	}

	dirPath, absErr := filepath.Abs(config.DirPath)
	if absErr != nil {
		return nil, errors.NewInvalidPathError(errors.InvalidPathConfig{
			Err:  absErr,
			Path: config.DirPath,
		})
	}

	return &Directory{dirPath: dirPath}, nil
}
