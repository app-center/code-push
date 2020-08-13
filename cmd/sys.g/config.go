package main

import (
	"errors"
	"fmt"
	"strings"
)

type serveConfig struct {
	ConfigFilePath string
	Debug          bool
	Port           int

	AddrCodePushD string
	AddrSessionD  string

	RootUserName string
	RootUserPwd  string
}

func (c *serveConfig) Validate() error {
	var errs []string

	if c.Port == 0 {
		errs = append(errs, "Invalid Port")
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
