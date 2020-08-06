package main

import (
	"errors"
	"fmt"
	"github.com/funnyecho/code-push/pkg/fs"
	"strings"
)

type serveConfig struct {
	ConfigFilePath string
	Debug          bool
	Port           int
	BoltPath       string

	AliOssEndpoint     string
	AliOssAccessKeyId  string
	AliOssAccessSecret string
}

func (c *serveConfig) validate() error {
	var errs []string

	if c.Port == 0 {
		errs = append(errs, "Invalid Port")
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
