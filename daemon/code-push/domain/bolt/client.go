package bolt

import (
	"fmt"
	"github.com/funnyecho/code-push/daemon/code-push/domain"
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
	"time"
)

type Client struct {
	Path string

	// Services
	branchService  BranchService
	envService     EnvService
	versionService VersionService

	db *bbolt.DB
}

func NewClient() *Client {
	c := &Client{}

	c.branchService.client = c
	c.envService.client = c
	c.versionService.client = c

	return c
}

func (c *Client) Open() error {
	if len(c.Path) == 0 {
		return fmt.Errorf("no database path provided")
	}

	// Open database file.
	db, err := bbolt.Open(c.Path, 0666, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return errors.Wrapf(err, "open bolt database failed, path: %s", c.Path)
	}
	c.db = db

	// Initialize top-level buckets.
	tx, err := c.db.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin writable tx failed while opening bolt database")
	}
	defer tx.Rollback()

	if _, err := tx.CreateBucketIfNotExists(bucketBranch); err != nil {
		return errors.Wrap(err, "create Branch bucket failed")
	}

	if _, err := tx.CreateBucketIfNotExists(bucketEnv); err != nil {
		return errors.Wrap(err, "create Env bucket failed")
	}

	if _, err := tx.CreateBucketIfNotExists(bucketEnvVersions); err != nil {
		return errors.Wrap(err, "create EnvVersions bucket failed")
	}

	return tx.Commit()
}

func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *Client) BranchService() *BranchService {
	return &c.branchService
}

func (c *Client) EnvService() *EnvService {
	return &c.envService
}

func (c *Client) VersionService() *VersionService {
	return &c.versionService
}

func (c *Client) DomainService() domain.Service {
	return &struct {
		*BranchService
		*EnvService
		*VersionService
	}{
		BranchService:  c.BranchService(),
		EnvService:     c.EnvService(),
		VersionService: c.VersionService(),
	}
}
