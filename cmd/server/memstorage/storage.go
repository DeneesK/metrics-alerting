package memstorage

import (
	"fmt"
	"strconv"
)

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

func (storage *MemStorage) Value(t, name string) string {
	switch t {
	case counterMetric:
		v, ok := storage.Counter[name]
		if !ok {
			return ""
		}
		return fmt.Sprintf("%d", v)
	case gaugeMetric:
		v, ok := storage.Gauge[name]
		if !ok {
			return ""
		}
		return fmt.Sprintf("%f", v)
	}
	return ""
}

func (storage *MemStorage) Metrics() string {
	result := ""
	for k, v := range storage.Counter {
		result += fmt.Sprintf("[%s]: %d\n", k, v)
	}
	for k, v := range storage.Gauge {
		result += fmt.Sprintf("[%s]: %f\n", k, v)
	}
	return result
}
