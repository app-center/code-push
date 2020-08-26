package noopfactory

import "github.com/funnyecho/code-push/pkg/metrics"

var NullGauge metrics.Gauge = &nullGauge{}

type nullGauge struct{}

func (n nullGauge) With(labelValues ...string) metrics.Gauge {
	return n
}

func (n nullGauge) Set(value float64) {}

func (n nullGauge) Add(delta float64) {}
