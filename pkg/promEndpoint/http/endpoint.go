package prometheus_http

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func Handle(serveMux *http.ServeMux) {
	serveMux.Handle("/metrics", promhttp.Handler())
}
