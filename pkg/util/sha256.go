package util

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/funnyecho/code-push/pkg/fs"
	"github.com/pkg/errors"
	"io"
	"os"
)

func EncodeFileSha256(path string) (string, error) {
	file, fileErr := fs.File(fs.FileConfig{FilePath: path})

	if fileErr != nil {
		return "", errors.Wrapf(fileErr, "invalid file, path:%s", path)
	}

	f, err := os.Open(file.Path())
	if err != nil {
		return "", err
	}

	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
