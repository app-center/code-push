package grpckit_metrics

import (
	"github.com/funnyecho/code-push/pkg/metrics"
)

func NewMetrics() *Metrics {
	m := &Metrics{}

	metrics.MustInitD(m)
	return m
}

type Metrics struct {
	GRPCRequestDuration metrics.Histogram `metric:"grpc_request_duration_seconds" labels:"method,errCode" help:"Duration of api requesting"`
	GRPCRequestSucceed  metrics.Counter   `metric:"grpc_request_total" tags:"result=ok" labels:"method,errCode" help:"grpc requesting counter"`
	GRPCRequestFailed   metrics.Counter   `metric:"grpc_request_total" tags:"result=err" labels:"method,errCode" help:"grpc requesting counter"`
}

func (m *Metrics) ObserveGRPCRequestDuration(seconds float64, method, errCode string) {
	m.GRPCRequestDuration.With(method, errCode).Observe(seconds)
}

func (m *Metrics) IncGRPCRequestSucceed(method string) {
	m.GRPCRequestSucceed.With(method, "").Inc()
}

func (m *Metrics) IncGRPCRequestFailed(method, errCode string) {
	m.GRPCRequestFailed.With(method, errCode).Inc()
}
