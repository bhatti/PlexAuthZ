package benchmark

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"sort"
	"time"
)

// SystemUsage tracks cpu/memory/disk usage.
type SystemUsage struct {
	CPUUtilization float64
	MemoryTotalMB  uint64
	MemoryUsedMB   uint64
	MemoryPercent  float64
	DiskTotalGB    uint64
	DiskUsedGB     uint64
	DiskPercent    float64
}

// String helper.
func (u *SystemUsage) String() string {
	return fmt.Sprintf(
		"[CPU: %f, Total Memory: %v MB, Used Memory: %v MB, Memory Util: %.2f%%, Total Disk: %v GB, Used Disk: %v GB, Disk Util: %.2f%%]",
		u.CPUUtilization,
		u.MemoryTotalMB,
		u.MemoryUsedMB,
		u.MemoryPercent,
		u.DiskTotalGB,
		u.DiskUsedGB,
		u.DiskPercent,
	)
}

// Metric tracks start/end time with latency, error and usage.
type Metric struct {
	Started time.Time
	Latency time.Duration
	Error   error
	Usage   SystemUsage
}

func newMetric() *Metric {
	return &Metric{
		Started: time.Now(),
	}
}

// String helper.
func (m *Metric) String() string {
	return fmt.Sprintf(
		"Started: %v\nLatency: %v\nError: %v\nSystem Usage: %s\n",
		m.Started,
		m.Latency,
		m.Error,
		m.Usage.String(),
	)
}

// Populate initializes cpu/memory/disk usage.
func (u *SystemUsage) Populate() {
	cpuPer, _ := cpu.Percent(0, false)
	u.CPUUtilization = cpuPer[0]
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		u.MemoryTotalMB = memInfo.Total / 1024 / 1024
		u.MemoryUsedMB = memInfo.Used / 1024 / 1024
		u.MemoryPercent = memInfo.UsedPercent
	}
	diskInfo, err := disk.Usage(".")
	if err == nil {
		u.DiskTotalGB = diskInfo.Total / 1024 / 1024 / 1024
		u.DiskUsedGB = diskInfo.Used / 1024 / 1024 / 1024
		u.DiskPercent = diskInfo.UsedPercent
	}
}

func (m *Metric) finish(err error) {
	m.Latency = time.Since(m.Started)
	m.Error = err
	m.Usage.Populate()
}

func calculateMetricPercentiles(metrics []Metric, percentiles []float64) (res map[float64]*Metric) {
	res = make(map[float64]*Metric)
	var memoryTotal []uint64
	var memoryUsed []uint64
	var diskTotal []uint64
	var diskUsed []uint64
	var latencies []time.Duration
	for _, metric := range metrics {
		memoryTotal = append(memoryTotal, metric.Usage.MemoryTotalMB)
		memoryUsed = append(memoryUsed, metric.Usage.MemoryUsedMB)
		diskTotal = append(diskTotal, metric.Usage.DiskTotalGB)
		diskUsed = append(diskUsed, metric.Usage.DiskUsedGB)
		latencies = append(latencies, metric.Latency)
	}
	for _, percentile := range percentiles {
		metric := &Metric{}
		for i := 0; i < len(metrics); i++ {
			metric.Error = metrics[i].Error
		}
		metric.Usage.MemoryTotalMB = calculateIntPercentile(memoryTotal, percentile)
		metric.Usage.MemoryUsedMB = calculateIntPercentile(memoryUsed, percentile)
		if metric.Usage.MemoryTotalMB > 0 {
			metric.Usage.MemoryPercent = float64(metric.Usage.MemoryUsedMB * 100 / metric.Usage.MemoryTotalMB)
		}
		metric.Usage.DiskTotalGB = calculateIntPercentile(diskTotal, percentile)
		metric.Usage.DiskUsedGB = calculateIntPercentile(diskUsed, percentile)
		if metric.Usage.DiskTotalGB > 0 {
			metric.Usage.DiskPercent = float64(metric.Usage.DiskUsedGB * 100 / metric.Usage.DiskTotalGB)
		}
		metric.Latency = calculateLatencyPercentile(latencies, percentile)
		res[percentile] = metric
	}
	return
}

func calculateLatencyPercentile(latencies []time.Duration, percentile float64) time.Duration {
	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})
	index := int(float64(len(latencies)) * percentile / 100)
	if index < len(latencies) {
		return latencies[index]
	}
	return 0
}

func calculateIntPercentile(arr []uint64, percentile float64) uint64 {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})
	index := int(float64(len(arr)) * percentile / 100)
	if index < len(arr) {
		return arr[index]
	}
	return 0
}
