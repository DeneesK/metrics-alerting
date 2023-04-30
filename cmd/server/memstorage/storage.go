package memstorage

import "strconv"

type MemStorage struct {
	Gauge   map[string]float32
	Counter map[string]int
}

const (
	counterMetric string = "counter"
	gaugeMetric   string = "gauge"
)

func NewMemStorage() MemStorage {
	return MemStorage{Gauge: make(map[string]float32), Counter: make(map[string]int)}
}

func (storage *MemStorage) StoreMetrics(t, name, value string) {
	switch t {
	case counterMetric:
		v, _ := strconv.Atoi(value)
		storage.Counter[name] += v
	case gaugeMetric:
		v, _ := strconv.ParseFloat(value, 32)
		storage.Gauge[name] = float32(v)
	}
}
