package noopfactory

import "github.com/funnyecho/code-push/pkg/metrics"

var NullCounter metrics.Counter = &nullCounter{}

type nullCounter struct{}

func (n nullCounter) With(labelValues ...string) metrics.Counter {
	return n
}

func (n nullCounter) Inc()              {}
func (n nullCounter) Add(delta float64) {}
