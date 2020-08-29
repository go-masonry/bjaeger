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
	// TraceIDKey key used in log
	TraceIDKey = "traceId"
	// SpanIDKey key used in log
	SpanIDKey = "spanId"
	// ParentSpanIDKey key used in log
	ParentSpanIDKey = "parentSpanId"
	// SampledKey key used in log
	SampledKey = "sampled"
)

// TraceInfoContextExtractorFxOption is a preconfigured fx.Option that will allow adding trace info to log entry
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

// TraceInfoExtractorFromContext helper function to extract trace info from context
func TraceInfoExtractorFromContext(ctx context.Context) map[string]interface{} {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		if jaegerContext, ok := span.Context().(jaeger.SpanContext); ok {
			return extractFromSpanContext(jaegerContext)
		}
	}
	return nil
}

func extractFromSpanContext(ctx jaeger.SpanContext) map[string]interface{} {
	var output = make(map[string]interface{}, 4)
	output[SpanIDKey] = ctx.SpanID().String()
	output[ParentSpanIDKey] = ctx.ParentID().String()
	output[SampledKey] = ctx.IsSampled()
	if traceID := ctx.TraceID(); traceID.IsValid() {
		output[TraceIDKey] = traceID.String()
	}
	return output
}
