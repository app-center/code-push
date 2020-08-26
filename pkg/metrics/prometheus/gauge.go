package promfactory

import (
	"github.com/funnyecho/code-push/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type gauge struct {
	gv  *prometheus.GaugeVec
	lbs []string
}

func (c *gauge) With(labelValues ...string) metrics.Gauge {
	return &gauge{gv: c.gv, lbs: append(c.lbs, labelValues...)}
}

func (c *gauge) Set(value float64) {
	c.gv.WithLabelValues(c.lbs...).Set(value)
}

func (c *gauge) Add(delta float64) {
	c.gv.WithLabelValues(c.lbs...).Add(delta)
}
