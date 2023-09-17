package metrics

import (
	"github.com/bhatti/PlexAuthZ/internal/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	dto "github.com/prometheus/client_model/go"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

// Registry keeps track of metrics
type Registry struct {
	registry   *prometheus.Registry
	counters   map[string]*prometheus.CounterVec
	gauges     map[string]*prometheus.GaugeVec
	histograms map[string]*prometheus.Histogram
	lock       sync.RWMutex
}

// New metrics constructor
func New() *Registry {
	registry := &Registry{
		registry:   prometheus.NewRegistry(),
		counters:   make(map[string]*prometheus.CounterVec),
		gauges:     make(map[string]*prometheus.GaugeVec),
		histograms: make(map[string]*prometheus.Histogram),
	}
	if err := registry.registry.Register(collectors.NewGoCollector()); err != nil {
		logrus.WithFields(logrus.Fields{
			"Component": "Metrics",
			"Error":     err,
		}).Warn("failed to register GO collector")
	}
	if err := registry.registry.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
		logrus.WithFields(logrus.Fields{
			"Component": "Metrics",
			"Error":     err,
		}).Warn("failed to register process collector")
	}
	return registry
}

// Incr metric
func (r *Registry) Incr(id string, args ...string) {
	opts := utils.ArrayToMap(args...)
	r.lock.Lock()
	defer r.lock.Unlock()
	counter := r.counters[id]
	if counter == nil {
		keys := make([]string, len(opts))
		i := 0
		for k := range opts {
			keys[i] = k
			i++
		}
		counter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: id + "_total",
			},
			keys,
		)
		if err := r.registry.Register(counter); err != nil {
			logrus.WithFields(logrus.Fields{
				"Component": "MetricsRegistry",
				"ID":        id,
				"Opts":      opts,
				"Error":     err,
			}).
				Error("failed to register counter")
			return
		}
		r.counters[id] = counter
	}
	labels := make([]string, len(opts))
	i := 0
	for _, v := range opts {
		labels[i] = v
		i++
	}
	counter.WithLabelValues(labels...).Inc()
}

// Elapsed metric.
func (r *Registry) Elapsed(id string, args ...string) func() {
	start := time.Now().Unix()
	return func() {
		elapsed := time.Now().Unix() - start
		r.Duration(id, float64(elapsed), args...)
	}
}

// Duration metric for latency.
func (r *Registry) Duration(id string, value float64, args ...string) {
	opts := utils.ArrayToMap(args...)
	r.lock.Lock()
	defer r.lock.Unlock()
	histogram := r.histograms[id]
	if histogram == nil {
		keys := make([]string, len(opts))
		i := 0
		for k := range opts {
			keys[i] = k
			i++
		}
		durations := prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    id + "_duration_seconds",
			Help:    id,
			Buckets: prometheus.ExponentialBuckets(0.1, 1.5, 5),
			//Buckets:                     prometheus.LinearBuckets(normMean-5*normDomain, .5*normDomain, 20),
			NativeHistogramBucketFactor: 1.1,
		})
		histogram = &durations
		if err := r.registry.Register(durations); err != nil {
			logrus.WithFields(logrus.Fields{
				"Component": "MetricsRegistry",
				"ID":        id,
				"Opts":      opts,
				"Error":     err,
			}).
				Error("failed to register histogram")
			return
		}
		r.histograms[id] = histogram
	}
	(*histogram).(prometheus.ExemplarObserver).ObserveWithExemplar(value, opts)
}

// Set gauge
func (r *Registry) Set(id string, val float64, args ...string) {
	opts := utils.ArrayToMap(args...)
	r.lock.Lock()
	defer r.lock.Unlock()
	gauge := r.gauges[id]
	if gauge == nil {
		keys := make([]string, len(opts))
		i := 0
		for k := range opts {
			keys[i] = k
			i++
		}
		gauge = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: id,
			},
			keys,
		)
		r.registry.MustRegister(gauge)
		r.gauges[id] = gauge
	}
	labels := make([]string, len(opts))
	i := 0
	for _, v := range opts {
		labels[i] = v
		i++
	}
	gauge.WithLabelValues(labels...).Set(val)
}

// Summary of metrics.
func (r *Registry) Summary() map[string]float64 {
	r.lock.Lock()
	defer r.lock.Unlock()
	res := make(map[string]float64)
	metrics, _ := r.registry.Gather()
	for _, metric := range metrics {
		for _, m := range metric.Metric {
			if metric.Name == nil || metric.Type == nil {
				continue
			}
			if *metric.Type == dto.MetricType_HISTOGRAM {
				if m.Histogram.SampleSum != nil {
					res[*metric.Name] = *m.Histogram.SampleSum
					res[strings.ReplaceAll(*metric.Name, "duration_seconds", "counts")] = float64(*m.Histogram.SampleCount)
				}
			} else if *metric.Type == dto.MetricType_COUNTER { // dto.MetricType_GAUGE
				if m.Counter.Value != nil {
					res[*metric.Name] = *m.Counter.Value
				}
			}
		}
	}
	return res
}
