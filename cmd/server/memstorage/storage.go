package memstorage

import "github.com/DeneesK/metrics-alerting/cmd/server/urlparser"

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

func SaveMetrics(u string, storage *MemStorage) error {
	m, err := urlparser.ParseMetricURL(u)
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
