package bolt

import (
	"github.com/funnyecho/code-push/daemon"
	"github.com/funnyecho/code-push/daemon/domain/bolt/internal"
	"github.com/pkg/errors"
	"time"
)

type VersionService struct {
	client *Client
}

func (s *VersionService) Version(envId, appVersion []byte) (*daemon.Version, error) {
	tx, err := s.client.db.Begin(false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin tx")
	}
	defer tx.Rollback()

	bucket := tx.Bucket(bucketEnvVersions).Bucket(envId)
	if bucket == nil {
		return nil, nil
	}

	var ver daemon.Version
	if v := bucket.Get(appVersion); v == nil {
		return nil, nil
	} else if err := internal.UnmarshalVersion(v, &ver); err != nil {
		return nil, err
	}

	return &ver, nil
}

func (s *VersionService) VersionsWithEnvId(envId []byte) (daemon.VersionList, error) {
	tx, err := s.client.db.Begin(false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin tx")
	}
	defer tx.Rollback()

	bucket := tx.Bucket(bucketEnvVersions).Bucket(envId)
	if bucket == nil {
		return nil, nil
	}

	var list daemon.VersionList

	c := bucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var ver daemon.Version
		if err := internal.UnmarshalVersion(v, &ver); err != nil {
			return nil, err
		}

		list = append(list, &ver)
	}

	return list, nil
}

func (s *VersionService) CreateVersion(version *daemon.Version) error {
	if len(version.EnvId) == 0 ||
		len(version.AppVersion) == 0 {
		return daemon.ErrParamsInvalid
	}

	if !s.client.EnvService().IsEnvAvailable([]byte(version.EnvId)) {
		return daemon.ErrEnvNotFound
	}

	tx, err := s.client.db.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin writable tx failed")
	}
	defer tx.Rollback()

	b, bucketErr := tx.Bucket(bucketEnvVersions).CreateBucketIfNotExists([]byte(version.EnvId))
	if bucketErr != nil {
		return errors.Wrap(bucketErr, "create env version bucket failed")
	}

	if v := b.Get([]byte(version.AppVersion)); v != nil {
		return errors.WithMessagef(
			daemon.ErrVersionExisted,
			"envId: %s, appVersion: %s",
			version.EnvId,
			version.AppVersion,
		)
	}

	version.CreateTime = time.Now()
	if v, err := internal.MarshalVersion(version); err != nil {
		return err
	} else if err := b.Put([]byte(version.AppVersion), v); err != nil {
		return errors.Wrap(err, "put version to tx failed")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit tx failed")
	}

	return nil
}

func (s *VersionService) IsVersionAvailable(envId, appVersion []byte) (bool, error) {
	version, err := s.Version(envId, appVersion)

	if err != nil {
		return false, err
	}

	return version != nil, nil
}
