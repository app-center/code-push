package fs

import (
	"github.com/funnyecho/code-push/pkg/fs/errors"
	"os"
	"path"
	"path/filepath"
)

type File struct {
	filePath string
}

func (f *File) Path() string {
	return f.filePath
}

func (f *File) Size() (size int64, err error) {
	fs, err := os.Stat(f.filePath)

	if err != nil {
		return 0, err
	}

	return fs.Size(), nil
}

func (f *File) Extension() string {
	return path.Ext(f.filePath)
}

func (f *File) CheckNotExist() bool {
	_, err := os.Stat(f.filePath)
	return os.IsNotExist(err)
}

func (f *File) CheckPermissionDenied() bool {
	_, err := os.Stat(f.filePath)
	return os.IsPermission(err)
}

func (f *File) DirPath() string {
	return filepath.Dir(f.filePath)
}

func (f *File) Directory() (*Directory, error) {
	return NewDirectory(DirectoryConfig{DirPath: f.DirPath()})
}

func (f *File) EnsurePath() error {
	dir, _ := f.Directory()
	return dir.EnsurePath()
}

func (f *File) Open(flag int, perm os.FileMode) (*os.File, error) {
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

func (f *File) MustOpen() (*os.File, error) {
	return f.Open(os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
}

func (f *File) Move(dstPath string) (*File, error) {
	dstFile, dstErr := NewFile(FileConfig{FilePath: dstPath})
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

func (f *File) Delete() error {
	if f.CheckNotExist() {
		return nil
	}

	return os.Remove(f.Path())
}

type FileConfig struct {
	FilePath string
}

func NewFile(config FileConfig) (f *File, err error) {
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

	return &File{filePath}, nil
}
