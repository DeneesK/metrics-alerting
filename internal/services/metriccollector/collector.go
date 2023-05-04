package metriccollector

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/fatih/structs"
)

var stats = []string{"Alloc", "BuckHashSys", "Frees", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc"}

type Collector interface {
	StartCollect()
	GetMetrics() map[string]float64
	GetRandomValue() float64
	GetPollCount() int
}

type Metrics struct {
	stats        runtime.MemStats
	pollCount    int
	randomValue  float64
	pollInterval time.Duration
}

func NewCollector(pollInterval int) Metrics {
	return Metrics{pollInterval: time.Duration(pollInterval)}
}

func (ms *Metrics) pollMetrics() {
	runtime.ReadMemStats(&ms.stats)
	ms.randomValue = float64(rand.Float32())
	ms.pollCount += 1
	time.Sleep(ms.pollInterval * time.Second)
}

func (ms *Metrics) StartCollect() {
	for {
		ms.pollMetrics()
	}
}

func (ms *Metrics) GetMetrics() map[string]float64 {
	m := structs.Map(ms.stats)
	mapstats := make(map[string]float64)
	for _, stat := range stats {
		mapstats[stat] = float64(m[stat].(uint64))
	}
	mapstats["GCCPUFraction"] = m["GCCPUFraction"].(float64)
	mapstats["NumForcedGC"] = float64(m["NumForcedGC"].(uint32))
	mapstats["NumGC"] = float64(m["NumGC"].(uint32))
	return mapstats
}

func (ms *Metrics) GetRandomValue() float64 {
	return float64(ms.randomValue)
}

func (ms *Metrics) GetPollCount() int {
	return int(ms.pollCount)
}
