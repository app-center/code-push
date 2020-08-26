package metrics

// Options defines the information associated with a metric
type Options struct {
	Name   string
	Tags   map[string]string
	Labels []string
	Help   string
}

// HistogramOptions defines the information associated with a metric
type HistogramOptions struct {
	Name   string
	Tags   map[string]string
	Labels []string
	Help   string
}

// Factory creates new metrics
type Factory interface {
	Counter(metric Options) Counter
	Gauge(metric Options) Gauge
	Histogram(metric HistogramOptions) Histogram
}

var NullFactory Factory = &nullFactory{}
var DefaultFactory Factory = NullFactory

func SetDefaultFactory(factory Factory) {
	DefaultFactory = factory
}

type nullFactory struct{}

func (f *nullFactory) Counter(options Options) Counter {
	return NullCounter
}

func (f *nullFactory) Gauge(options Options) Gauge {
	return NullGauge
}

func (f *nullFactory) Histogram(options HistogramOptions) Histogram {
	return NullSummary
}
