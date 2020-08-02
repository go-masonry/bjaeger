package bjaeger

import (
	"context"
	"io"
	"sync"

	"github.com/go-masonry/mortar/interfaces/trace"
	"github.com/opentracing/opentracing-go"
)

type tracerWrapper struct {
	sync.Mutex
	cfg               *jaegerConfig
	tracer            opentracing.Tracer
	closer            io.Closer
	closed, connected bool
}

func newWrapper(cfg *jaegerConfig) trace.OpenTracer {
	return &tracerWrapper{
		cfg:    cfg,
		tracer: opentracing.NoopTracer{},
	}
}

// tracer wrapper impl
func (w *tracerWrapper) Connect(ctx context.Context) error {
	w.Lock()
	defer w.Unlock()
	if !w.connected {
		var connectionErrorChannel = make(chan error)
		go func() {
			tracer, closer, err := w.cfg.conf.NewTracer(w.cfg.options...)
			if err == nil {
				w.tracer = tracer
				w.closer = closer
				w.connected = true
			}
			connectionErrorChannel <- err
		}()
		select {
		case err := <-connectionErrorChannel:
			return err
		case <-ctx.Done():
			w.connected = false
			w.tracer = opentracing.NoopTracer{}
			return ctx.Err()
		}
	}
	return nil
}

func (w *tracerWrapper) Close(ctx context.Context) error {
	w.Lock()
	defer func() {
		w.closed = true
		w.Unlock()
	}()
	if w.closed || !w.connected || w.closer == nil {
		return nil
	}
	doneChannel := make(chan error)

	go func(ch chan error) {
		ch <- w.closer.Close()
	}(doneChannel)

	select {
	case err := <-doneChannel:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (w *tracerWrapper) Tracer() opentracing.Tracer {
	return w
}

func (w *tracerWrapper) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return w.tracer.StartSpan(operationName, opts...)
}

func (w *tracerWrapper) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return w.tracer.Inject(sm, format, carrier)
}

func (w *tracerWrapper) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return w.tracer.Extract(format, carrier)
}
