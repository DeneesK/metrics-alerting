package metric

import (
	"math/rand"
	"runtime"
	"time"
)

type gauge float32
type counter int

type MetricStats struct {
	Stats        runtime.MemStats
	PollCount    counter
	RandomValue  gauge
	pollInterval time.Duration
}

func NewMetricStats() MetricStats {
	return MetricStats{pollInterval: 2}
}

func (ms *MetricStats) PollStats() {
	runtime.ReadMemStats(&ms.Stats)
	ms.RandomValue = gauge(rand.Float32())
	ms.PollCount += 1
	time.Sleep(ms.pollInterval * time.Second)
}

func (ms *MetricStats) StartCollect() {
	for {
		ms.PollStats()
	}
}
