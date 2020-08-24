package client

import "github.com/funnyecho/code-push/pkg/metrics"

type Metrics struct {
	RequestDuration metrics.Histogram `meta:""`
}
