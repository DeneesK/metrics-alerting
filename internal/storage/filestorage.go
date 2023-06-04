package storage

import (
	"encoding/json"
	"os"
)

type allMetrics struct {
	Gauge   map[string]float64 `json:"gauge"`
	Counter map[string]int64   `json:"counter"`
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func newProducer(filename string) (*producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &producer{file: file, encoder: json.NewEncoder(file)}, nil
}

func (p *producer) writeMetrics(m *allMetrics) error {
	return p.encoder.Encode(&m)
}

func (p *producer) close() error {
	return p.file.Close()
}

type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func newConsumer(filename string) (*consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &consumer{file: file, decoder: json.NewDecoder(file)}, nil
}

func (c *consumer) readMetrics() (*allMetrics, error) {
	var data allMetrics
	err := c.decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *consumer) close() error {
	return c.file.Close()
}

func storeToFile(path string, m *MemStorage) error {
	p, err := newProducer(path)
	if err != nil {
		return err
	}
	defer p.close()
	var metrics allMetrics

	metrics.Counter = m.GetCounterMetrics()
	metrics.Gauge = m.GetGaugeMetrics()

	return p.writeMetrics(&metrics)
}

func loadFromFile(path string, m *MemStorage) error {
	c, err := newConsumer(path)
	if err != nil {
		return err
	}
	defer c.close()
	metrics, err := c.readMetrics()
	if err != nil {
		return err
	}
	m.setMetrics(metrics)
	return nil
}
