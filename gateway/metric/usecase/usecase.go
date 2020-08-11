package usecase

import (
	"github.com/funnyecho/code-push/pkg/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func New() UseCase {
	uc := &useCase{}

	uc.initMetrics()

	return uc
}

type useCase struct {
	requestDurationMetric metrics.Histogram
}

func (uc *useCase) initMetrics() {
	uc.requestDurationMetric = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "code_push",
		Name: "request_duration_seconds",
		Help: "Request duration in seconds",
	}, []string{"svr_type", "svr_name", "interface", "path", "success"})
}
