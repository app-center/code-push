package noopfactory

import "github.com/funnyecho/code-push/pkg/metrics"

var NullFactory metrics.Factory = &factory{}

type factory struct{}

func (f *factory) Counter(options metrics.Options) metrics.Counter {
	return NullCounter
}

func (f *factory) Gauge(options metrics.Options) metrics.Gauge {
	return NullGauge
}

func (f *factory) Histogram(options metrics.HistogramOptions) metrics.Histogram {
	return NullSummary
}
