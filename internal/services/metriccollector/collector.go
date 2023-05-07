package metriccollector

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	memStatsLen = 15
	cpuStatsLen = 1
)

type Metrics struct {
	mx           sync.Mutex
	stats        runtime.MemStats
	pollCount    int64
	randomValue  float64
	pollInterval time.Duration
}

type RuntimeMetrics struct {
	memMetrics map[string]uint64
	cpuMetrics map[string]float64
}

func (rm *RuntimeMetrics) GetMemMetrics() map[string]uint64 {
	memMetricsCopy := make(map[string]uint64, len(rm.memMetrics))
	for k, v := range rm.memMetrics {
		memMetricsCopy[k] = v
	}
	return memMetricsCopy
}

func (rm *RuntimeMetrics) GetCPUMetrics() map[string]float64 {
	cpuMetricCopy := make(map[string]float64, len(rm.cpuMetrics))
	for k, v := range rm.cpuMetrics {
		cpuMetricCopy[k] = v
	}
	return cpuMetricCopy
}

func NewCollector(pollInterval int) Metrics {
	return Metrics{pollInterval: time.Duration(pollInterval) * time.Second}
}

func (ms *Metrics) pollMetrics() {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	runtime.ReadMemStats(&ms.stats)
	ms.randomValue = float64(rand.Float32())
}

func (ms *Metrics) StartCollect() {
	for {
		ms.pollMetrics()
		ms.pollCount += 1
		time.Sleep(ms.pollInterval)
	}
}

func (ms *Metrics) GetRuntimeMetrics() RuntimeMetrics {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	memStats := make(map[string]uint64, memStatsLen)
	memStats["Alloc"] = ms.stats.Alloc
	memStats["BuckHashSys"] = ms.stats.BuckHashSys
	memStats["Frees"] = ms.stats.Frees
	memStats["GCSys"] = ms.stats.GCSys
	memStats["HeapAlloc"] = ms.stats.HeapAlloc
	memStats["HeapIdle"] = ms.stats.HeapIdle
	memStats["HeapInuse"] = ms.stats.HeapInuse
	memStats["HeapObjects"] = ms.stats.HeapObjects
	memStats["HeapReleased"] = ms.stats.HeapReleased
	memStats["HeapSys"] = ms.stats.HeapSys
	memStats["Lookups"] = ms.stats.Lookups
	memStats["MCacheInuse"] = ms.stats.MCacheInuse
	memStats["MCacheSys"] = ms.stats.MCacheSys
	memStats["MSpanInuse"] = ms.stats.MSpanInuse
	memStats["MSpanSys"] = ms.stats.MSpanSys
	memStats["NextGC"] = ms.stats.NextGC
	memStats["Mallocs"] = ms.stats.Mallocs
	memStats["OtherSys"] = ms.stats.OtherSys
	memStats["PauseTotalNs"] = ms.stats.PauseTotalNs
	memStats["StackInuse"] = ms.stats.StackInuse
	memStats["StackSys"] = ms.stats.StackSys
	memStats["Sys"] = ms.stats.Sys
	memStats["TotalAlloc"] = ms.stats.TotalAlloc
	memStats["NumForcedGC"] = uint64(ms.stats.NumForcedGC)
	memStats["NumGC"] = uint64(ms.stats.NumGC)

	cpuStats := make(map[string]float64, cpuStatsLen)
	cpuStats["GCCPUFraction"] = ms.stats.GCCPUFraction

	return RuntimeMetrics{memMetrics: memStats, cpuMetrics: cpuStats}

}

func (ms *Metrics) GetRandomValue() float64 {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	return ms.randomValue
}

func (ms *Metrics) GetPollCount() int64 {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	return ms.pollCount
}
