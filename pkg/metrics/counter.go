package metrics

// Counter describes a metric that accumulates values monotonically.
// An example of a counter is the number of received HTTP requests.
type Counter interface {
	With(labelValues ...string) Counter
	Inc()
	Add(delta float64)
}

var NullCounter Counter = &nullCounter{}

type nullCounter struct{}

func (n nullCounter) With(labelValues ...string) Counter {
	return n
}

func (n nullCounter) Inc()              {}
func (n nullCounter) Add(delta float64) {}
