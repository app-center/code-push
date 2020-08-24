package opentracing

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func StartTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, spanCtx := opentracing.StartSpanFromContext(c.Request.Context(), c.Request.URL.Path)
		defer span.Finish()

		span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		c.Request = c.Request.Clone(spanCtx)

		c.Next()
	}
}
