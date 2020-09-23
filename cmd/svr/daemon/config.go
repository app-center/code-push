package main

import (
	"errors"
	"fmt"
	"github.com/funnyecho/code-push/pkg/fs"
	"strings"
)

type serveConfig struct {
	ConfigFilePath string `flag:"config" value:"config/daemon.yml" usage:"alternative config file path"`
	Debug          bool   `flag:"debug" value:"false" usage:"run in debug mode"`

	PortGrpc       int    `flag:"port_grpc" usage:"port for grpc server listen to"`
	PortHttp       int    `flag:"port_http" usage:"port for http server listen to"`

	BoltPath       string `flag:"bbolt_path" value:"storage/daemon/bbolt.db" usage:"path of bbolt storage file"`

	AliOssEndpoint     string `flag:"alioss_endpoint" usage:"endpoint of ali-oss"`
	AliOssBucket       string `flag:"alioss_bucket" usage:"bucket of ali-oss"`
	AliOssAccessKeyId  string `flag:"alioss_access_key_id" usage:"access key id of ali-oss"`
	AliOssAccessSecret string `flag:"alioss_access_secret" usage:"access secret of ali-oss"`
}

func (c *serveConfig) Validate() error {
	var errs []string

	if c.PortGrpc == 0 {
		errs = append(errs, "Invalid Grpc Port")
	}

	if c.PortHttp == 0 {
		errs = append(errs, "Invalid Http Port")
	}

	if c.BoltPath == "" {
		errs = append(errs, "BoltPath required")
	} else {
		boltFile, boltFileErr := fs.File(fs.FileConfig{
			FilePath: c.BoltPath,
		})
		if boltFileErr != nil {
			errs = append(errs, boltFileErr.Error())
		} else {
			if dirErr := boltFile.EnsurePath(); dirErr != nil {
				errs = append(errs, dirErr.Error())
			}
		}
	}

	if c.AliOssEndpoint == "" {
		errs = append(errs, "AliOssEndpoint required")
	}

	if c.AliOssBucket == "" {
		errs = append(errs, "AliOssBucket required")
	}

	if c.AliOssAccessKeyId == "" {
		errs = append(errs, "AliOssAccessKeyId required")
	}

	if c.AliOssAccessSecret == "" {
		errs = append(errs, "AliOssAccessSecret required")
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(fmt.Sprintf("FA_CONFIG_SERVE:\n\t%s", strings.Join(errs[:], "\n\t")))
}
