package metrics

// Gauge describes a metric that takes specific values over time.
// An example of a gauge is the current depth of a job queue.
type Gauge interface {
	With(labelValues ...string) Gauge
	Set(value float64)
	Add(delta float64)
}

var NullGauge Gauge = &nullGauge{}

type nullGauge struct{}

func (n nullGauge) With(labelValues ...string) Gauge {
	return n
}

func (n nullGauge) Set(value float64) {}

func (n nullGauge) Add(delta float64) {}
