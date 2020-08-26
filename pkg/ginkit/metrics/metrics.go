package ginkit_metrics

import (
	"github.com/funnyecho/code-push/pkg/metrics"
	"strconv"
)

func NewMetrics() *Metrics {
	m := &Metrics{}

	metrics.MustInitD(m)
	return m
}

type Metrics struct {
	HttpRequestDuration metrics.Histogram `metric:"http_request_duration_seconds" labels:"method,path,status_code" help:"Duration of api requesting"`
	HttpRequestSucceed  metrics.Counter   `metric:"http_request_total" tags:"result=ok" labels:"method,path,status_code"`
	HttpRequestFailed   metrics.Counter   `metric:"http_request_total" tags:"result=err" labels:"method,path,status_code"`
}

func (m *Metrics) ObserveHttpRequestDuration(seconds float64, method, path string, statusCode int) {
	m.HttpRequestDuration.
		With(method, path, strconv.Itoa(statusCode)).
		Observe(seconds)
}

func (m *Metrics) IncHttpRequestSucceed(method, path string, statusCode int) {
	m.HttpRequestSucceed.
		With(method, path, strconv.Itoa(statusCode)).
		Inc()
}

func (m *Metrics) IncHttpRequestFailed(method, path string, statusCode int) {
	m.HttpRequestFailed.
		With(method, path, strconv.Itoa(statusCode)).
		Inc()
}
