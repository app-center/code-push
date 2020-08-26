package metrics

// Histogram describes a metric that takes repeated observations of the same
// kind of thing, and produces a statistical summary of those observations,
// typically expressed as quantiles or buckets. An example of a histogram is
// HTTP request latencies.
type Histogram interface {
	With(labelValues ...string) Histogram
	Observe(value float64)
}

var NullSummary Histogram = &nullSummary{}

type nullSummary struct{}

func (s nullSummary) With(labelValues ...string) Histogram {
	return s
}

func (s nullSummary) Observe(value float64) {}
