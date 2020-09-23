package bolt

import (
	"github.com/funnyecho/code-push/daemon"
	"github.com/funnyecho/code-push/daemon/domain/bolt/internal"
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
	"time"
)

type EnvService struct {
	client *Client
}

func (s *EnvService) Env(envId []byte) (*daemon.Env, error) {
	tx, err := s.client.db.Begin(false)
	if err != nil {
		return nil, errors.Wrap(err, "begin read-only tx failed")
	}
	defer tx.Rollback()

	var e daemon.Env
	if v := tx.Bucket(bucketEnv).Get(envId); v == nil {
		return nil, nil
	} else if err := internal.UnmarshalEnv(v, &e); err != nil {
		return nil, err
	}

	return &e, nil
}

func (s *EnvService) CreateEnv(env *daemon.Env) error {
	if len(env.ID) == 0 ||
		len(env.Name) == 0 ||
		len(env.EncToken) == 0 ||
		len(env.BranchId) == 0 {
		return daemon.ErrParamsInvalid
	}

	if !s.client.BranchService().IsBranchAvailable([]byte(env.BranchId)) {
		return daemon.ErrBranchNotFound
	}

	tx, err := s.client.db.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin writable tx failed")
	}
	defer tx.Rollback()

	b := tx.Bucket(bucketEnv)
	if v := b.Get([]byte(env.ID)); v != nil {
		return errors.WithMessagef(
			daemon.ErrEnvExisted,
			"envId: %s",
			env.ID,
		)
	}

	env.CreateTime = time.Now()
	if v, err := internal.MarshalEnv(env); err != nil {
		return err
	} else if err := b.Put([]byte(env.ID), v); err != nil {
		return errors.Wrap(err, "put env to tx failed")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit tx failed")
	}

	return nil
}

func (s *EnvService) DeleteEnv(envId []byte) error {
	if len(envId) == 0 {
		return errors.WithMessage(daemon.ErrParamsInvalid, "envId required")
	}

	tx, err := s.client.db.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin write tx failed")
	}
	defer tx.Rollback()

	b := tx.Bucket(bucketEnv)
	if err := b.Delete(envId); err != nil {
		return errors.WithMessagef(
			err,
			"delete env failed, envId: %s",
			envId,
		)
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit tx failed")
	}

	return nil
}

func (s *EnvService) IsEnvAvailable(envId []byte) bool {
	env, err := s.Env(envId)

	if err != nil || env == nil {
		return false
	}

	return s.client.BranchService().IsBranchAvailable([]byte(env.BranchId))
}

func (s *EnvService) IsEnvNameExisted(branchId, envName []byte) (bool, error) {
	tx, err := s.client.db.Begin(false)
	if err != nil {
		return false, errors.Wrap(err, "failed to begin tx")
	}
	defer tx.Rollback()

	existed := false
	err = s.client.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketEnv)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var env daemon.Env
			if err := internal.UnmarshalEnv(v, &env); err != nil {
				return err
			} else if env.BranchId == string(branchId) && env.Name == string(envName) {
				existed = true
			}
		}

		return nil
	})

	return existed, err
}
