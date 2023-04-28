package memstorage

import "github.com/DeneesK/metrics-alerting/cmd/server/urlparser"

type Repository interface {
	StoreMetrics(string) error
}

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

func (storage *MemStorage) StoreMetrics(str string) error {
	m, err := urlparser.ParseMetricURL(str)
	if err != nil {
		return err
	}
	switch m.MetricType {
	case counterMetric:
		storage.Counter[m.MetricName] += m.CounterValue
	case gaugeMetric:
		storage.Gauge[m.MetricName] = m.CaugeValue
	}
	return nil
}
