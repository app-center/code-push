package main

import (
	"errors"
	"fmt"
	"strings"
)

type serveConfig struct {
	ConfigFilePath string `flag:"config" value:"config/serve.yml" usage:"alternative config file path"`
	Debug          bool   `flag:"sys_g.debug" value:"false" usage:"run in debug mode"`
	Port           int    `flag:"sys_g.port_http" usage:"port for http server listen to"`

	AddrCodePushD string `flag:"addr_code_push_d" usage:"address of code-push.d"`
	AddrSessionD  string `flag:"addr_session_d" usage:"address of session.d"`

	RootUserName string `flag:"sys_g.root_user_name" usage:"root user name"`
	RootUserPwd  string `flag:"sys_g.root_user_pwd" usage:"root user password"`
}

func (c *serveConfig) Validate() error {
	var errs []string

	if c.Port == 0 {
		errs = append(errs, "Invalid PortGrpc")
	}

	if c.AddrCodePushD == "" {
		errs = append(errs, "Invalid address of code-push.d")
	}

	if c.AddrSessionD == "" {
		errs = append(errs, "Invalid address of session.d")
	}

	if len(c.RootUserName) < 6 {
		errs = append(errs, "length of root use name shall be larger than 6")
	}

	if len(c.RootUserPwd) < 8 {
		errs = append(errs, "length of root use pwd shall be larger than 8")
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(fmt.Sprintf("FA_CONFIG_SERVE:\n\t%s", strings.Join(errs[:], "\n\t")))
}
