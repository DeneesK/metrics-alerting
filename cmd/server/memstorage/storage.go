package memstorage

import (
	"github.com/DeneesK/metrics-alerting/cmd/server/urlparser"
)

type MemStorage struct {
	Gauge   map[string]float32
	Counter map[string]int
}

func NewMemStorage() MemStorage {
	return MemStorage{Gauge: make(map[string]float32), Counter: make(map[string]int)}
}

func SaveMetrics(m urlparser.MetricPayload, storage *MemStorage) {
	if m.MetricType == "counter" {
		storage.Counter[m.MetricName] += m.CounterValue
		return
	}
	storage.Gauge[m.MetricName] = m.CaugeValue
}
