package promfactory

import (
	"github.com/funnyecho/code-push/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type counter struct {
	cv  *prometheus.CounterVec
	lbs []string
}

func (c *counter) With(labelValues ...string) metrics.Counter {
	return &counter{cv: c.cv, lbs: append(c.lbs, labelValues...)}
}

func (c *counter) Inc() {
	c.cv.WithLabelValues(c.lbs...).Inc()
}

func (c *counter) Add(delta float64) {
	c.cv.WithLabelValues(c.lbs...).Add(delta)
}
