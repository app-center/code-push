package main

import (
	"errors"
	"fmt"
	"strings"
)

type serveConfig struct {
	ConfigFilePath string `flag:"config" value:"config/serve.yml" usage:"alternative config file path"`
	Debug          bool   `flag:"session_d.debug" value:"false" usage:"run in debug mode"`
	PortGrpc       int    `flag:"session_d.port_grpc" usage:"port for grpc server listen to"`
	PortHttp       int    `flag:"session_d.port_http" usage:"port for http server listen to"`
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
