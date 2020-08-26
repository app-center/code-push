package tracing

import (
	"fmt"
	"github.com/funnyecho/code-push/pkg/log"
	"github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"io"
	"time"
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

	registry := stdprometheus.NewRegistry()
	pusher := push.New("http://localhost:9091", "code_push_opentracing").Gatherer(registry)

	tracer, closer, err = tracerConfig.NewTracer(
		config.Logger(&tracerLogger{logger: logger}),
		config.Metrics(prometheus.New(
			prometheus.WithRegisterer(registry),
		)),
	)

	{
		ticker := time.NewTicker(15 * time.Second)
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					pusher.Add()
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()

		closer = &tracerCloser{
			pusherQuit: quit,
			oriClose:   closer,
		}
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
	pusherQuit chan struct{}
	oriClose   io.Closer
}

func (t *tracerCloser) Close() error {
	t.pusherQuit <- struct{}{}
	return t.oriClose.Close()
}
