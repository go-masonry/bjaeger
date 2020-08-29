package bjaeger

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber/jaeger-client-go/config"
)

func TestBuilder(t *testing.T) {
	builder, err := Builder().SetServiceName("myname").AddOptions(BricksLoggerOption(nil)).
		SetCustomConfig(&config.Configuration{}).
		Build()
	assert.NoError(t, err)
	wrapper, ok := builder.(*tracerWrapper)
	assert.True(t, ok, "wrong type")
	assert.Len(t, wrapper.cfg.options, 1)
	assert.NotNil(t, wrapper.cfg.conf)
}

func TestDefaultBuilder(t *testing.T) {
	_, err := Builder().Build()
	assert.NoError(t, err)
}

func TestConnect(t *testing.T) {
	os.Setenv("JAEGER_DISABLE", "true")
	tracer, err := Builder().SetServiceName("name").Build()
	assert.NoError(t, err)
	err = tracer.Connect(context.Background())
	assert.NoError(t, err)
}

func TestConnectFails(t *testing.T) {
	tracer, err := Builder().Build()
	assert.NoError(t, err)
	err = tracer.Connect(context.Background())
	assert.Error(t, err)
}

func TestConnectContextDone(t *testing.T) {
	doneCtx, cancel := context.WithCancel(context.Background())
	cancel()
	os.Setenv("JAEGER_DISABLE", "true")
	tracer, err := Builder().SetServiceName("name").Build()
	assert.NoError(t, err)
	err = tracer.Connect(doneCtx)
	assert.Error(t, err)
}

func TestClose(t *testing.T) {
	os.Setenv("JAEGER_DISABLE", "true")
	tracer, err := Builder().SetServiceName("name").Build()
	assert.NoError(t, err)
	err = tracer.Connect(context.Background())
	assert.NoError(t, err)
	err = tracer.Close(context.Background())
	assert.NoError(t, err)
}

func TestCloseWithDoneContext(t *testing.T) {
	os.Setenv("JAEGER_DISABLE", "true")
	tracer, err := Builder().SetServiceName("name").Build()
	assert.NoError(t, err)
	err = tracer.Connect(context.Background())
	assert.NoError(t, err)
	doneCtx, cancel := context.WithCancel(context.Background())
	cancel()
	err = tracer.Close(doneCtx)
	assert.Error(t, err)
}

func TestStartSpan(t *testing.T) {
	// will work since by default it's NoopTracer
	tracer := newWrapper(new(jaegerConfig))
	span := tracer.Tracer().StartSpan("name")
	assert.NotNil(t, span)
}

func TestInject(t *testing.T) {
	// will work since by default it's NoopTracer
	tracer := newWrapper(new(jaegerConfig))
	err := tracer.Tracer().Inject(nil, nil, nil)
	assert.NoError(t, err)
}

func TestExtract(t *testing.T) {
	// will work since by default it's NoopTracer
	tracer := newWrapper(new(jaegerConfig))
	_, err := tracer.Tracer().Extract(nil, nil)
	// span not found
	assert.Error(t, err)
}
