package promfactory

import (
	"github.com/funnyecho/code-push/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type summary struct {
	sv  *prometheus.SummaryVec
	lbs []string
}

func (s *summary) With(labelValues ...string) metrics.Histogram {
	return &summary{
		sv:  s.sv,
		lbs: append(s.lbs, labelValues...),
	}
}

func (s *summary) Observe(value float64) {
	s.sv.WithLabelValues(s.lbs...).Observe(value)
}
