package bolt

import (
	"github.com/funnyecho/code-push/daemon/filer"
	"github.com/funnyecho/code-push/daemon/filer/domain/bolt/internal"
	"github.com/pkg/errors"
	"time"
)

type FileService struct {
	client *Client
}

func (s *FileService) File(fileKey string) (*filer.File, error) {
	if fileKey == "" {
		return nil, filer.ErrInvalidFileKey
	}

	tx, err := s.client.db.Begin(false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin tx")
	}
	defer tx.Rollback()

	var f filer.File
	if v := tx.Bucket(bucketFile).Get([]byte(fileKey)); v == nil {
		return nil, nil
	} else if err := internal.UnmarshalFile(v, &f); err != nil {
		return nil, err
	}

	return &f, nil
}

func (s *FileService) InsertFile(file *filer.File) error {
	if file == nil {
		return filer.ErrParamsInvalid
	}

	if file.Key == "" {
		return filer.ErrInvalidFileKey
	}

	if file.Value == "" {
		return filer.ErrInvalidFileValue
	}

	tx, err := s.client.db.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin writable tx failed")
	}
	defer tx.Rollback()

	b := tx.Bucket(bucketFile)
	if v := b.Get([]byte(file.Key)); v != nil {
		return errors.WithMessagef(
			filer.ErrFileKeyExisted,
			"fileKey: %s",
			file.Key,
		)
	}

	file.CreateTime = time.Now()

	if v, err := internal.MarshalFile(file); err != nil {
		return err
	} else if err := b.Put([]byte(file.Key), v); err != nil {
		return errors.Wrap(err, "put file to tx failed")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit tx failed")
	}

	return nil
}

func (s *FileService) IsFileKeyExisted(fileKey string) bool {
	f, err := s.File(fileKey)

	return err == nil && f != nil
}
