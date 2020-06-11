package bolt

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain"
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

	return c
}

func (c *Client) Open() error {
	// Open database file.
	db, err := bbolt.Open(c.Path, 0666, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	c.db = db

	// Initialize top-level buckets.
	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.CreateBucketIfNotExists([]byte("Dials")); err != nil {
		return err
	}

	return tx.Commit()
}

func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *Client) BranchService() domain.IBranchService {
	return &c.branchService
}

func (c *Client) EnvService() domain.IEnvService {
	return &c.envService
}

func (c *Client) VersionService() domain.IVersionService {
	return &c.versionService
}
