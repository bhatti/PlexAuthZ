package benchmark

import (
	"fmt"
	"sync"
	"time"
)

// FuncToBenchmark for Benchmark function
type FuncToBenchmark func() error

// Request for Benchmark function
type Request struct {
	TPS         int
	Duration    time.Duration
	Percentiles []float64
}

// Response for Benchmark function
type Response struct {
	StartedAt            time.Time
	BeforeBenchmarkUsage SystemUsage
	AfterBenchmarkUsage  SystemUsage
	ErrorMetrics         map[float64]*Metric
	OKMetrics            map[float64]*Metric
	ErrorsByType         map[string]int
	TPS                  int
	ActualTPS            int
}

func (r *Response) String() string {
	return fmt.Sprintf(
		"Error Metrics:\n%v\nOK Metrics:\n%v\nError Types:\n%v\nTPS %d, Actual TPS %d:\nBefore Usage: %s\nAfter Usage: %s\n",
		r.ErrorMetrics,
		r.OKMetrics,
		r.ErrorsByType,
		r.TPS,
		r.ActualTPS,
		r.BeforeBenchmarkUsage.String(),
		r.AfterBenchmarkUsage.String(),
	)
}
func (r *Response) populate(req Request, okMetrics []Metric, errorMetrics []Metric) {
	errorsByTape := make(map[string]int)
	for _, metric := range errorMetrics {
		if metric.Error != nil {
			err := metric.Error.Error()
			if len(err) > 20 {
				err = err[0:20]
			}
			errorsByTape[err] = errorsByTape[err] + 1
		}
	}
	r.ActualTPS = int(float64(len(okMetrics)+len(errorMetrics)) / req.Duration.Seconds())
	r.OKMetrics = calculateMetricPercentiles(okMetrics, req.Percentiles)
	r.ErrorMetrics = calculateMetricPercentiles(errorMetrics, req.Percentiles)
	r.ErrorsByType = errorsByTape
	r.AfterBenchmarkUsage.Populate()
}

// Benchmark function
func Benchmark(fn FuncToBenchmark, req Request) (res Response) {
	var errorMetrics []Metric
	var okMetrics []Metric
	// for synchronization
	wg := &sync.WaitGroup{}
	var lock sync.Mutex

	// for receiving metrics
	metricCh := make(chan *Metric)

	// calculating tick time
	throttle := time.Tick(time.Second / time.Duration(req.TPS))

	res.BeforeBenchmarkUsage.Populate()
	res.TPS = req.TPS
	res.StartedAt = time.Now()

	go func() {
		for metric := range metricCh {
			lock.Lock()
			if metric.Error != nil {
				errorMetrics = append(errorMetrics, *metric)
			} else {
				okMetrics = append(okMetrics, *metric)
			}
			lock.Unlock()
		}
	}()

	for startTime := time.Now(); time.Since(startTime) <= req.Duration; {
		<-throttle
		wg.Add(1)
		go func() {
			defer wg.Done()
			metric := newMetric()
			err := fn()
			metric.finish(err)
			metricCh <- metric
		}()
	}
	wg.Wait()
	close(metricCh)
	// send back response
	lock.Lock()
	defer lock.Unlock()
	time.Sleep(time.Second) // wait a second to return system to normal before after snapshot
	res.populate(req, okMetrics, errorMetrics)
	return
}
