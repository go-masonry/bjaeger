package bjaeger

import (
	"context"

	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/providers/groups"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/fx"
)

const (
	TraceIDKey      = "traceId"
	SpanIDKey       = "spanId"
	ParentSpanIDKey = "parentSpanId"
	SampledKey      = "sampled"
)

func TraceInfoContextExtractorFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group: groups.LoggerContextExtractors,
			Target: func() log.ContextExtractor {
				return TraceInfoExtractorFromContext
			},
		},
	)
}

func TraceInfoExtractorFromContext(ctx context.Context) map[string]interface{} {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		if jaegerContext, ok := span.Context().(jaeger.SpanContext); ok {
			var output = make(map[string]interface{}, 4)
			output[SpanIDKey] = jaegerContext.SpanID().String()
			output[ParentSpanIDKey] = jaegerContext.ParentID().String()
			output[SampledKey] = jaegerContext.IsSampled()
			if traceID := jaegerContext.TraceID(); traceID.IsValid() {
				output[TraceIDKey] = traceID.String()
			}
			return output
		}
	}
	return nil
}
