package main

import (
	"errors"
	"fmt"
	"strings"
)

type serveConfig struct {
	ConfigFilePath string `flag:"config" value:"config/serve.yml" usage:"alternative config file path"`
	Debug          bool   `flag:"portal_g.debug" value:"false" usage:"run in debug mode"`
	Port           int    `flag:"portal_g.port_http" usage:"port for http server listen to"`

	AddrCodePushD string `flag:"addr_code_push_d" usage:"address of code-push.d"`
	AddrFilerD    string `flag:"addr_filer_d" usage:"address of filer.d"`
	AddrSessionD  string `flag:"addr_session_d" usage:"address of session.d"`
}

func (c *serveConfig) Validate() error {
	var errs []string

	if c.Port == 0 {
		errs = append(errs, "Invalid PortGrpc")
	}

	if c.AddrCodePushD == "" {
		errs = append(errs, "Invalid address of code-push.d")
	}

	if c.AddrFilerD == "" {
		errs = append(errs, "Invalid address of filer.d")
	}

	if c.AddrSessionD == "" {
		errs = append(errs, "Invalid address of session.d")
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(fmt.Sprintf("FA_CONFIG_SERVE:\n\t%s", strings.Join(errs[:], "\n\t")))
}
