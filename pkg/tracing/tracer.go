package tracing

import (
	"fmt"
	"github.com/funnyecho/code-push/pkg/log"
	"github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"io"
)

func InitTracer(serviceName string, logger log.Logger) (tracer opentracing.Tracer, closer io.Closer, err error) {
	tracerConfig, tracerConfigErr := config.FromEnv()
	if tracerConfigErr != nil {
		return nil, nil, tracerConfigErr
	}

	tracerConfig.ServiceName = serviceName
	tracerConfig.Sampler.Type = "const"
	tracerConfig.Sampler.Param = 1

	tracerConfig.Reporter.LogSpans = true
	tracerConfig.RPCMetrics = true

	tracer, closer, err = tracerConfig.NewTracer(
		config.Logger(&tracerLogger{logger: logger}),
		config.Metrics(prometheus.New(prometheus.WithRegisterer(stdprometheus.DefaultRegisterer))),
	)

	closer = &tracerCloser{
		oriClose: closer,
	}

	return
}

type tracerLogger struct {
	logger log.Logger
}

func (t *tracerLogger) Error(msg string) {
	t.logger.Info(msg, "error", true)
}

func (t *tracerLogger) Infof(msg string, args ...interface{}) {
	t.logger.Debug(fmt.Sprintf(msg, args...))
}

type tracerCloser struct {
	oriClose io.Closer
}

func (t *tracerCloser) Close() error {
	return t.oriClose.Close()
}
