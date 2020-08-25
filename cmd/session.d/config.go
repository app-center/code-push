package main

import (
	"errors"
	"fmt"
	"strings"
)

type serveConfig struct {
	ConfigFilePath string
	Debug          bool
	PortGrpc       int
	PortHttp       int
}

func (c *serveConfig) Validate() error {
	var errs []string

	if c.PortGrpc == 0 {
		errs = append(errs, "Invalid Grpc Port")
	}

	if c.PortHttp == 0 {
		errs = append(errs, "Invalid Http Port")
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(fmt.Sprintf("FA_CONFIG_SERVE:\n\t%s", strings.Join(errs[:], "\n\t")))
}
