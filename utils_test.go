package bjaeger

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uber/jaeger-client-go"
)

func TestTraceInfoExtractorFromContext(t *testing.T) {
	ctx := jaeger.NewSpanContext(jaeger.TraceID{High: 255, Low: 255}, jaeger.SpanID(1), jaeger.SpanID(1), true, nil)
	values := extractFromSpanContext(ctx)
	assert.Contains(t, values, "traceId")
	assert.Contains(t, values, "spanId")
	assert.Contains(t, values, "parentSpanId")
	assert.Contains(t, values, "sampled")
}
