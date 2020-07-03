package bolt

import (
	"github.com/funnyecho/code-push/daemon/filer"
	"github.com/funnyecho/code-push/daemon/filer/domain"
	"github.com/funnyecho/code-push/daemon/filer/domain/bolt/internal"
	"github.com/pkg/errors"
	"time"
)

type SchemeService struct {
	client *Client
}

func (s *SchemeService) AliOssScheme() (*domain.AliOssScheme, error) {
	tx, err := s.client.db.Begin(false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin tx")
	}
	defer tx.Rollback()

	var a domain.AliOssScheme
	if v := tx.Bucket(bucketScheme).Get(keyAliOssScheme); v == nil {
		return nil, nil
	} else if err := internal.UnmarshalAliOssScheme(v, &a); err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *SchemeService) UpdateAliOssScheme(scheme *domain.AliOssScheme) error {
	fetchedScheme, fetchErr := s.AliOssScheme()
	if fetchErr != nil {
		return fetchErr
	}

	var schemeToUpdate domain.AliOssScheme
	if fetchedScheme != nil {
		schemeToUpdate = *fetchedScheme
	}

	if scheme.Endpoint != nil {
		schemeToUpdate.Endpoint = scheme.Endpoint
	}

	if scheme.AccessKeyId != nil {
		schemeToUpdate.AccessKeyId = scheme.AccessKeyId
	}

	if scheme.AccessKeySecret != nil {
		schemeToUpdate.AccessKeySecret = scheme.AccessKeySecret
	}

	if schemeToUpdate.Endpoint == nil {
		return filer.ErrInvalidAliOssEndpoint
	}

	if schemeToUpdate.AccessKeyId == nil {
		return filer.ErrInvalidAliOssAccessKeyId
	}

	if schemeToUpdate.AccessKeySecret == nil {
		return filer.ErrInvalidAliOssAccessKeySecret
	}

	schemeToUpdate.UpdateTime = time.Now()

	tx, err := s.client.db.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin writable tx failed")
	}
	defer tx.Rollback()

	if v, err := internal.MarshalAliOssScheme(&schemeToUpdate); err != nil {
		return err
	} else if err := tx.Bucket(bucketScheme).Put(keyAliOssScheme, v); err != nil {
		return errors.Wrap(err, "put ali oss scheme to tx failed")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit tx failed")
	}

	return nil
}
