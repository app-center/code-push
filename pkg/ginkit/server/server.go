package ginkit_server

import (
	ginkit_metrics "github.com/funnyecho/code-push/pkg/ginkit/metrics"
	ginkit_middleware_logger "github.com/funnyecho/code-push/pkg/ginkit/middleware/logger"
	ginkit_middleware_metrics "github.com/funnyecho/code-push/pkg/ginkit/middleware/metrics"
	ginkit_middleware_opentracing "github.com/funnyecho/code-push/pkg/ginkit/middleware/opentracing"
	"github.com/funnyecho/code-push/pkg/log"
	prometheus_gin "github.com/funnyecho/code-push/pkg/prom-endpoint/gin"
	"github.com/gin-gonic/gin"
)

type Options func(*options)

func New(opts ...Options) *gin.Engine {
	opt := &options{}

	for _, fn := range opts {
		fn(opt)
	}

	if opt.debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	if !opt.disablePrometheusEndpoint {
		prometheus_gin.Init(r)
	}

	if !opt.disableOpentracing {
		useOpentracing := ginkit_middleware_opentracing.UseOpentracing()
		r.Use(useOpentracing)
	}

	if opt.logger != nil {
		r.Use(ginkit_middleware_logger.UseLogger(opt.logger))
	}

	if !opt.disableMetrics {
		metrics := ginkit_metrics.NewMetrics()

		// global metrics middleware has covered no-method response
		//r.NoMethod(func(context *gin.Context) {
		//	metrics.IncHttpRequestFailed(context.Request.Method, context.Request.URL.Path, context.Writer.Status())
		//})

		// global metrics middleware has covered no-route response
		//r.NoRoute(func(context *gin.Context) {
		//	metrics.IncHttpRequestFailed(context.Request.Method, context.Request.URL.Path, context.Writer.Status())
		//})

		r.Use(ginkit_middleware_metrics.UseMetrics(metrics))
	}

	return r
}

func WithLogger(logger log.Logger) Options {
	return func(o *options) {
		o.logger = logger
	}
}

func WithDebugMode(debug bool) Options {
	return func(o *options) {
		o.debugMode = debug
	}
}

func WithDisablePrometheusEndpoint() Options {
	return func(o *options) {
		o.disablePrometheusEndpoint = true
	}
}

func WithDisableMetrics() Options {
	return func(o *options) {
		o.disableMetrics = true
	}
}

func WithDisableOpentracing() Options {
	return func(o *options) {
		o.disableOpentracing = true
	}
}

type options struct {
	debugMode                 bool
	disablePrometheusEndpoint bool
	disableMetrics            bool
	disableOpentracing        bool
	logger                    log.Logger
}
