package bolt

import (
	"fmt"
	"github.com/funnyecho/code-push/daemon/filer/domain"
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
	"time"
)

func NewClient() *Client {
	client := &Client{}

	client.fileService.client = client

	return client
}

type Client struct {
	Path string

	fileService FileService

	db *bbolt.DB
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

	if _, err := tx.CreateBucketIfNotExists(bucketFile); err != nil {
		return errors.Wrap(err, "create file bucket failed")
	}

	return tx.Commit()
}

func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *Client) FileService() domain.FileService {
	return &c.fileService
}

func (c *Client) DomainService() *domain.Service {
	return &domain.Service{
		FileService: c.FileService(),
	}
}
