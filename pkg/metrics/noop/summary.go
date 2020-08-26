package noopfactory

import "github.com/funnyecho/code-push/pkg/metrics"

var NullSummary metrics.Histogram = &nullSummary{}

type nullSummary struct{}

func (s nullSummary) With(labelValues ...string) metrics.Histogram {
	return s
}

func (s nullSummary) Observe(value float64) {}
