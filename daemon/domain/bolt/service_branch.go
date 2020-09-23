package bolt

import (
	"github.com/funnyecho/code-push/daemon"
	"github.com/funnyecho/code-push/daemon/domain/bolt/internal"
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
	"time"
)

type BranchService struct {
	client *Client
}

func (s *BranchService) Branch(branchId []byte) (*daemon.Branch, error) {
	tx, err := s.client.db.Begin(false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin tx")
	}
	defer tx.Rollback()

	var b daemon.Branch
	if v := tx.Bucket(bucketBranch).Get(branchId); v == nil {
		return nil, nil
	} else if err := internal.UnmarshalBranch(v, &b); err != nil {
		return nil, err
	}

	return &b, nil
}

func (s *BranchService) CreateBranch(branch *daemon.Branch) error {
	if len(branch.ID) == 0 ||
		len(branch.Name) == 0 ||
		len(branch.EncToken) == 0 {
		return daemon.ErrParamsInvalid
	}

	tx, err := s.client.db.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin writable tx failed")
	}
	defer tx.Rollback()

	b := tx.Bucket(bucketBranch)
	if v := b.Get([]byte(branch.ID)); v != nil {
		return errors.WithMessagef(
			daemon.ErrBranchExisted,
			"branchId: %s",
			branch.ID,
		)
	}

	branch.CreateTime = time.Now()

	if v, err := internal.MarshalBranch(branch); err != nil {
		return err
	} else if err := b.Put([]byte(branch.ID), v); err != nil {
		return errors.Wrap(err, "put branch to tx failed")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit tx failed")
	}

	return nil
}

func (s *BranchService) DeleteBranch(branchId []byte) error {
	if len(branchId) == 0 {
		return errors.WithMessage(daemon.ErrParamsInvalid, "branchId required")
	}

	tx, err := s.client.db.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin write tx failed")
	}
	defer tx.Rollback()

	b := tx.Bucket(bucketBranch)
	if err := b.Delete(branchId); err != nil {
		return errors.WithMessagef(
			err,
			"delete branch failed, branchId: %s",
			branchId,
		)
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit tx failed")
	}

	return nil
}

func (s *BranchService) IsBranchAvailable(branchId []byte) bool {
	branch, err := s.Branch(branchId)

	return err == nil && branch != nil
}

func (s *BranchService) IsBranchNameExisted(branchName []byte) (bool, error) {
	tx, err := s.client.db.Begin(false)
	if err != nil {
		return false, errors.Wrap(err, "failed to begin tx")
	}
	defer tx.Rollback()

	existed := false
	err = s.client.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketBranch)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var branch daemon.Branch
			if err := internal.UnmarshalBranch(v, &branch); err != nil {
				return err
			} else if branch.Name == string(branchName) {
				existed = true
			}
		}

		return nil
	})

	return existed, err
}
