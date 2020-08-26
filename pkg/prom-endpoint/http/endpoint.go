package prometheus_http

import (
	prom_endpoint "github.com/funnyecho/code-push/pkg/prom-endpoint"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func Handle(serveMux *http.ServeMux) {
	serveMux.Handle(prom_endpoint.Endpoint, promhttp.Handler())
}
