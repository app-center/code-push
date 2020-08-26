package client

import (
	"github.com/funnyecho/code-push/pkg/metrics"
	"net/http"
)

func NewMetrics(factory metrics.Factory) *Metrics {
	m := &Metrics{}

	metrics.MustInit(m, factory)
	return m
}

type Metrics struct {
	ApiRequestDuration metrics.Histogram `metric:"api_request_duration_seconds" labels:"method,path,status_code" help:"Duration of api requesting"`
	ApiRequestSucceed  metrics.Counter   `metric:"api_request_total" tags:"result=ok" labels:"method,path,status_code"`
	ApiRequestFailed   metrics.Counter   `metric:"api_request_total" tags:"result=err" labels:"method,path,status_code"`
}

func (m *Metrics) ObserveApiRequestDuration(seconds float64, request http.Request, response http.Response) {
	m.ApiRequestDuration.
		With("method", request.Method, "path", request.URL.Path, "status_code", string(response.StatusCode)).
		Observe(seconds)
}

func (m *Metrics) IncApiRequestSucceed(request http.Request, response http.Response) {
	m.ApiRequestSucceed.
		With("method", request.Method, "path", request.URL.Path, "status_code", string(response.StatusCode)).
		Inc()
}

func (m *Metrics) IncApiRequestFailed(request http.Request, response http.Response) {
	m.ApiRequestFailed.
		With("method", request.Method, "path", request.URL.Path, "status_code", string(response.StatusCode)).
		Inc()
}
