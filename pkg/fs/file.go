package fs

import (
	"github.com/funnyecho/code-push/pkg/fs/errors"
	"os"
	"path"
	"path/filepath"
)

type file struct {
	filePath string
}

func (f *file) Path() string {
	return f.filePath
}

func (f *file) Size() (size int64, err error) {
	fs, err := os.Stat(f.filePath)

	if err != nil {
		return 0, err
	}

	return fs.Size(), nil
}

func (f *file) Extension() string {
	return path.Ext(f.filePath)
}

func (f *file) CheckNotExist() bool {
	_, err := os.Stat(f.filePath)
	return os.IsNotExist(err)
}

func (f *file) CheckPermissionDenied() bool {
	_, err := os.Stat(f.filePath)
	return os.IsPermission(err)
}

func (f *file) DirPath() string {
	return filepath.Dir(f.filePath)
}

func (f *file) Directory() (*directory, error) {
	return Directory(DirectoryConfig{DirPath: f.DirPath()})
}

func (f *file) EnsurePath() error {
	dir, _ := f.Directory()
	return dir.EnsurePath()
}

func (f *file) Open(flag int, perm os.FileMode) (*os.File, error) {
	dir, dirErr := f.Directory()

	if dirErr != nil {
		return nil, dirErr
	}

	dirErr = dir.EnsurePath()
	if dirErr != nil {
		return nil, dirErr
	}

	if dir.CheckPermissionDenied() {
		return nil, errors.NewPermissionDeniedError(errors.PermissionDeniedConfig{
			Path: dir.Path(),
		})
	}

	file, err := os.OpenFile(f.filePath, flag, perm)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f *file) MustOpen() (*os.File, error) {
	return f.Open(os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
}

func (f *file) Move(dstPath string) (*file, error) {
	dstFile, dstErr := File(FileConfig{FilePath: dstPath})
	if dstErr != nil {
		return nil, dstErr
	}

	dstErr = dstFile.EnsurePath()
	if dstErr != nil {
		return nil, dstErr
	}

	moveErr := os.Rename(f.Path(), dstFile.Path())
	if moveErr != nil {
		return nil, moveErr
	}

	return dstFile, nil
}

func (f *file) Delete() error {
	if f.CheckNotExist() {
		return nil
	}

	return os.Remove(f.Path())
}

type FileConfig struct {
	FilePath string
}

func File(config FileConfig) (f *file, err error) {
	if len(config.FilePath) <= 0 {
		return nil, errors.NewInvalidPathError(errors.InvalidPathConfig{
			Path: config.FilePath,
		})
	}

	filePath, absErr := filepath.Abs(config.FilePath)
	if absErr != nil {
		return nil, errors.NewInvalidPathError(errors.InvalidPathConfig{
			Err:  absErr,
			Path: config.FilePath,
		})
	}

	return &file{filePath}, nil
}
