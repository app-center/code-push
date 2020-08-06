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
