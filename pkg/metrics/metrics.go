package metrics

import (
	metricsKit "github.com/go-kit/kit/metrics"
)

type Counter metricsKit.Counter
type Gauge metricsKit.Gauge
type Histogram metricsKit.Histogram
