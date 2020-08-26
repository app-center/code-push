package client

import (
	"github.com/funnyecho/code-push/pkg/metrics"
)

func NewMetrics() *Metrics {
	m := &Metrics{}

	metrics.MustInitD(m)
	return m
}

type Metrics struct {
}
