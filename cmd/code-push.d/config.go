package main

import (
	"errors"
	"fmt"
	"github.com/funnyecho/code-push/pkg/fs"
	"strings"
)

type serveConfig struct {
	ConfigFilePath string `flag:"config" value:"config/serve.yml" usage:"alternative config file path"`
	Debug          bool   `flag:"code_push_d.debug" value:"false" usage:"run in debug mode"`
	PortGrpc       int    `flag:"code_push_d.port_grpc" usage:"port for grpc server listen to"`
	PortHttp       int    `flag:"code_push_d.port_http" usage:"port for http server listen to"`
	BoltPath       string `flag:"code_push_d.bbolt_path" value:"storage/code-push.d/bbolt.db" usage:"path of bbolt storage file"`
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

	if len(errs) == 0 {
		return nil
	}

	return errors.New(fmt.Sprintf("FA_CONFIG_SERVE:\n\t%s", strings.Join(errs[:], "\n\t")))
}
