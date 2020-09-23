package main

import (
	"errors"
	"fmt"
	"strings"
)

type serveConfig struct {
	ConfigFilePath string `flag:"config" value:"config/gateway.yml" usage:"alternative config file path"`
	Debug          bool   `flag:"debug" value:"false" usage:"run in debug mode"`
	Port           int    `flag:"port_http" usage:"port for http server listen to"`

	AddrDaemon string `flag:"addr_daemon" usage:"address of code-push.daemon"`

	RootUserName string `flag:"root_user_name" usage:"root user name"`
	RootUserPwd  string `flag:"root_user_pwd" usage:"root user password"`
}

func (c *serveConfig) Validate() error {
	var errs []string

	if c.Port == 0 {
		errs = append(errs, "Invalid PortGrpc")
	}

	if c.AddrDaemon == "" {
		errs = append(errs, "Invalid address of code-push.daemon")
	}

	if len(c.RootUserName) < 6 {
		errs = append(errs, "length of root use name shall be larger than 6")
	}

	if len(c.RootUserPwd) < 6 {
		errs = append(errs, "length of root use pwd shall be larger than 6")
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(fmt.Sprintf("FA_CONFIG_SERVE:\n\t%s", strings.Join(errs[:], "\n\t")))
}
