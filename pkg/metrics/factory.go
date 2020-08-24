package metrics

// Options defines the information associated with a metric
type Options struct {
	Name string
	Tags map[string]string
	Help string
}

// HistogramOptions defines the information associated with a metric
type HistogramOptions struct {
	Name    string
	Tags    map[string]string
	Help    string
	Buckets []float64
}

// Factory creates new metrics
type Factory interface {
	Counter(metric Options) Counter
	Gauge(metric Options) Gauge
	Histogram(metric HistogramOptions) Histogram
}
