package promfactory

import (
	"fmt"
	"github.com/funnyecho/code-push/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strings"
)

func New() metrics.Factory {
	return &factory{}
}

type factory struct{}

func (f *factory) Counter(opt metrics.Options) (metric metrics.Counter) {
	metricName := opt.Name

	if !strings.HasSuffix(metricName, "_total") {
		metricName = fmt.Sprintf("%s_total", metricName)
	}

	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(prometheus.AlreadyRegisteredError); ok {
				cv, _ := e.ExistingCollector.(*prometheus.CounterVec)
				metric = &counter{
					cv: cv,
				}
			} else {
				panic(err)
			}
		}
	}()

	cv := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:        metricName,
			Help:        opt.Help,
			ConstLabels: opt.Tags,
		},
		opt.Labels,
	)

	return &counter{cv: cv}
}

func (f *factory) Gauge(opt metrics.Options) (metric metrics.Gauge) {
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(prometheus.AlreadyRegisteredError); ok {
				gv, _ := e.ExistingCollector.(*prometheus.GaugeVec)
				metric = &gauge{
					gv: gv,
				}
			} else {
				panic(err)
			}
		}
	}()

	gv := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        opt.Name,
			Help:        opt.Help,
			ConstLabels: opt.Tags,
		},
		opt.Labels,
	)

	return &gauge{gv: gv}
}

func (f *factory) Histogram(opt metrics.HistogramOptions) (metric metrics.Histogram) {
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(prometheus.AlreadyRegisteredError); ok {
				sv, _ := e.ExistingCollector.(*prometheus.SummaryVec)
				metric = &summary{
					sv: sv,
				}
			} else {
				panic(err)
			}
		}
	}()

	sv := promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:        opt.Name,
			Help:        opt.Help,
			ConstLabels: opt.Tags,
		},
		opt.Labels,
	)

	return &summary{sv: sv}
}
