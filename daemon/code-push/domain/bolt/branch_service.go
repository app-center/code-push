package bolt

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain"
	"github.com/funnyecho/code-push/daemon/code-push/domain/bolt/internal"
	"github.com/pkg/errors"
	"time"
)

var _ domain.IBranchService = &BranchService{}

type BranchService struct {
	client *Client
}

func (s *BranchService) Branch(branchId string) (*domain.Branch, error) {
	tx, err := s.client.db.Begin(false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin tx")
	}
	defer tx.Rollback()

	var b domain.Branch
	if v := tx.Bucket([]byte("Branches")).Get([]byte(branchId)); v == nil {
		return nil, nil
	} else if err := internal.UnmarshalBranch(v, &b); err != nil {
		return nil, err
	}

	return &b, nil
}

func (s *BranchService) CreateBranch(branch *domain.Branch) error {
	if len(branch.ID) == 0 ||
		len(branch.Name) == 0 ||
		len(branch.AuthHost) == 0 ||
		len(branch.EncToken) == 0 {
		return ErrBranchCreationParamsInvalid
	}

	tx, err := s.client.db.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin write tx failed")
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte("Branches"))
	if v := b.Get([]byte(branch.ID)); v != nil {
		return errors.WithMessagef(
			ErrBranchExists,
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

func (s *BranchService) DeleteBranch(branchId string) error {
	if len(branchId) == 0 {
		return errors.WithMessagef(ErrBranchCreationParamsInvalid, "branchId: %s", branchId)
	}

	tx, err := s.client.db.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin write tx failed")
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte("Branches"))
	if err := b.Delete([]byte(branchId)); err != nil {
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
