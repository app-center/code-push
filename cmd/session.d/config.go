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
}

func (c *serveConfig) Validate() error {
	var errs []string

	if c.Port == 0 {
		errs = append(errs, "Invalid Port")
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(fmt.Sprintf("FA_CONFIG_SERVE:\n\t%s", strings.Join(errs[:], "\n\t")))
}
