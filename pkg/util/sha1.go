package util

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/funnyecho/code-push/pkg/errors"
	"github.com/funnyecho/code-push/pkg/fs"
	"io"
	"os"
)

func EncodeFileSha1(path string) (string, error) {
	file, fileErr := fs.File(fs.FileConfig{FilePath: path})

	if fileErr != nil {
		return "", errors.Throw(errors.CtorConfig{
			Error: fileErr,
			Msg:   "invalid file",
			Meta:  errors.MetaFields{"path": path},
		})
	}

	f, err := os.Open(file.Path())
	if err != nil {
		return "", err
	}

	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
