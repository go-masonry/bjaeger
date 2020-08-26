package bjaeger

import (
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
	newMetric := m.metrics.WithTags(monitor.Tags{"namespace": scope.Name})
	if len(scope.Tags) > 0 {
		newMetric = newMetric.WithTags(scope.Tags)
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
	metric := ci.metrics.WithTags(ci.Tags)
	metric.Counter(ci.Name, "").Add(float64(i))
}

type timerImpl struct {
	metrics monitor.Metrics
	metrics.TimerOptions
}

func (ti *timerImpl) Record(d time.Duration) {
	metric := ti.metrics.WithTags(ti.Tags)
	metric.Timer(ti.Name, "").Record(d)
}

type gaugeImpl struct {
	metrics monitor.Metrics
	metrics.Options
}

func (gi *gaugeImpl) Update(i int64) {
	metric := gi.metrics.WithTags(gi.Tags)
	metric.Gauge(gi.Name, "").Set(float64(i))
}

type histogramImpl struct {
	metrics monitor.Metrics
	metrics.HistogramOptions
}

func (hi *histogramImpl) Record(f float64) {
	metric := hi.metrics.WithTags(hi.Tags)
	metric.Histogram(hi.Name, "", hi.HistogramOptions.Buckets).Record(f)
}
