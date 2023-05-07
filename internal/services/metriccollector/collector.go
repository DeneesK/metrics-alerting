package metriccollector

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const uint64StatsLen = 15

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
	statsUint64 := make(map[string]uint64, uint64StatsLen)
	statsUint64["Alloc"] = ms.stats.Alloc
	statsUint64["BuckHashSys"] = ms.stats.BuckHashSys
	statsUint64["Frees"] = ms.stats.Frees
	statsUint64["GCSys"] = ms.stats.GCSys
	statsUint64["HeapAlloc"] = ms.stats.HeapAlloc
	statsUint64["HeapIdle"] = ms.stats.HeapIdle
	statsUint64["HeapInuse"] = ms.stats.HeapInuse
	statsUint64["HeapObjects"] = ms.stats.HeapObjects
	statsUint64["HeapReleased"] = ms.stats.HeapReleased
	statsUint64["HeapSys"] = ms.stats.HeapSys
	statsUint64["Lookups"] = ms.stats.Lookups
	statsUint64["MCacheInuse"] = ms.stats.MCacheInuse
	statsUint64["MCacheSys"] = ms.stats.MCacheSys
	statsUint64["MSpanInuse"] = ms.stats.MSpanInuse
	statsUint64["MSpanSys"] = ms.stats.MSpanSys
	statsUint64["NextGC"] = ms.stats.NextGC
	statsUint64["Mallocs"] = ms.stats.Mallocs
	statsUint64["OtherSys"] = ms.stats.OtherSys
	statsUint64["PauseTotalNs"] = ms.stats.PauseTotalNs
	statsUint64["StackInuse"] = ms.stats.StackInuse
	statsUint64["StackSys"] = ms.stats.StackSys
	statsUint64["Sys"] = ms.stats.Sys
	statsUint64["TotalAlloc"] = ms.stats.TotalAlloc
	statsUint64["NumForcedGC"] = uint64(ms.stats.NumForcedGC)
	statsUint64["NumGC"] = uint64(ms.stats.NumGC)

	statsFloat64 := make(map[string]float64)
	statsFloat64["GCCPUFraction"] = ms.stats.GCCPUFraction

	return RuntimeMetrics{memMetrics: statsUint64, cpuMetrics: statsFloat64}

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
