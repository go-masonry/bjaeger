package bjaeger

import (
	"context"
	"time"

	"github.com/go-masonry/mortar/interfaces/monitor"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

func BricksMetricsOption(bricksMetrics monitor.Metrics) (opt config.Option) {
	if bricksMetrics != nil {
		factory := &metricsWrapper{metrics: bricksMetrics}
		return config.Metrics(factory) // This code can panic...
	}
	return func(*config.Options) {} // empty option
}

type metricsWrapper struct {
	metrics monitor.Metrics
}

func (m *metricsWrapper) Counter(options metrics.Options) metrics.Counter {
	return &counterImpl{metrics: m.metrics, Options: options}
}

func (m *metricsWrapper) Timer(options metrics.TimerOptions) metrics.Timer {
	return &timerImpl{metrics: m.metrics, TimerOptions: options}
}

func (m *metricsWrapper) Gauge(options metrics.Options) metrics.Gauge {
	return &gaugeImpl{metrics: m.metrics, Options: options}
}

func (m *metricsWrapper) Histogram(options metrics.HistogramOptions) metrics.Histogram {
	return &histogramImpl{metrics: m.metrics, HistogramOptions: options}
}

func (m *metricsWrapper) Namespace(scope metrics.NSOptions) metrics.Factory {
	// Best effort here
	newMetric := m.metrics.AddTag("namespace", scope.Name)
	for name, value := range scope.Tags {
		newMetric = newMetric.AddTag(name, value)
	}
	return &metricsWrapper{
		metrics: newMetric,
	}
}

type counterImpl struct {
	metrics monitor.Metrics
	metrics.Options
}

func (ci *counterImpl) Inc(i int64) {
	metric := ci.metrics
	for name, value := range ci.Tags {
		metric = metric.AddTag(name, value)
	}
	metric.Count(context.Background(), ci.Name, i)
}

type timerImpl struct {
	metrics monitor.Metrics
	metrics.TimerOptions
}

func (ti *timerImpl) Record(d time.Duration) {
	metric := ti.metrics
	for name, value := range ti.Tags {
		metric = metric.AddTag(name, value)
	}
	metric.Timing(context.Background(), ti.Name, d)
}

type gaugeImpl struct {
	metrics monitor.Metrics
	metrics.Options
}

func (gi *gaugeImpl) Update(i int64) {
	metric := gi.metrics
	for name, value := range gi.Tags {
		metric = metric.AddTag(name, value)
	}
	metric.Gauge(context.Background(), gi.Name, float64(i))
}

type histogramImpl struct {
	metrics monitor.Metrics
	metrics.HistogramOptions
}

func (hi *histogramImpl) Record(f float64) {
	metric := hi.metrics
	for name, value := range hi.Tags {
		metric = metric.AddTag(name, value)
	}
	metric.Histogram(context.Background(), hi.Name, f)
}
