package bolt

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain"
	"github.com/funnyecho/code-push/daemon/code-push/domain/bolt/internal"
	"github.com/pkg/errors"
	"time"
)

var _ domain.IEnvService = &EnvService{}

type EnvService struct {
	client *Client
}

func (s *EnvService) Env(envId string) (*domain.Env, error) {
	tx, err := s.client.db.Begin(false)
	if err != nil {
		return nil, errors.Wrap(err, "begin read-only tx failed")
	}
	defer tx.Rollback()

	var e domain.Env
	if v := tx.Bucket(bucketEnv).Get([]byte(envId)); v == nil {
		return nil, nil
	} else if err := internal.UnmarshalEnv(v, &e); err != nil {
		return nil, err
	}

	return &e, nil
}

func (s *EnvService) CreateEnv(env *domain.Env) error {
	if len(env.ID) == 0 ||
		len(env.Name) == 0 ||
		len(env.EncToken) == 0 ||
		len(env.BranchId) == 0 {
		return domain.ErrEnvCreationParamsInvalid
	}

	if !s.client.BranchService().IsBranchAvailable(env.BranchId) {
		return domain.ErrBranchNotFound
	}

	tx, err := s.client.db.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin writable tx failed")
	}
	defer tx.Rollback()

	b := tx.Bucket(bucketEnv)
	if v := b.Get([]byte(env.ID)); v != nil {
		return errors.WithMessagef(
			domain.ErrEnvExists,
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

func (s *EnvService) DeleteEnv(envId string) error {
	if len(envId) == 0 {
		return errors.WithMessage(domain.ErrParamsInvalid, "envId required")
	}

	tx, err := s.client.db.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin write tx failed")
	}
	defer tx.Rollback()

	b := tx.Bucket(bucketEnv)
	if err := b.Delete([]byte(envId)); err != nil {
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

func (s *EnvService) IsEnvAvailable(envId string) bool {
	env, err := s.Env(envId)

	if err != nil || env == nil {
		return false
	}

	return s.client.BranchService().IsBranchAvailable(env.BranchId)
}
