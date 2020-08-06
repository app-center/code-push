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

	PortCodePushD int
	PortFilerD    int
	PortSessionD  int
}

func (c *serveConfig) validate() error {
	var errs []string

	if c.Port == 0 {
		errs = append(errs, "Invalid Port")
	}

	if c.PortCodePushD == 0 {
		errs = append(errs, "Invalid port of code-push.d")
	}

	if c.PortFilerD == 0 {
		errs = append(errs, "Invalid port of filer.d")
	}

	if c.PortSessionD == 0 {
		errs = append(errs, "Invalid port of session.d")
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(fmt.Sprintf("FA_CONFIG_SERVE:\n\t%s", strings.Join(errs[:], "\n\t")))
}
