package bolt

import (
	"github.com/funnyecho/code-push/daemon"
	"github.com/funnyecho/code-push/daemon/domain/bolt/internal"
	"github.com/pkg/errors"
	"time"
)

type FileService struct {
	client *Client
}

func (s *FileService) File(fileKey string) (*daemon.File, error) {
	if fileKey == "" {
		return nil, daemon.ErrInvalidFileKey
	}

	tx, err := s.client.db.Begin(false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin tx")
	}
	defer tx.Rollback()

	var f daemon.File
	if v := tx.Bucket(bucketFile).Get([]byte(fileKey)); v == nil {
		return nil, nil
	} else if err := internal.UnmarshalFile(v, &f); err != nil {
		return nil, err
	}

	return &f, nil
}

func (s *FileService) InsertFile(file *daemon.File) error {
	if file == nil {
		return daemon.ErrParamsInvalid
	}

	if file.Key == "" {
		return daemon.ErrInvalidFileKey
	}

	if file.Value == "" {
		return daemon.ErrInvalidFileValue
	}

	tx, err := s.client.db.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin writable tx failed")
	}
	defer tx.Rollback()

	b := tx.Bucket(bucketFile)
	if v := b.Get([]byte(file.Key)); v != nil {
		return errors.WithMessagef(
			daemon.ErrFileKeyExisted,
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
