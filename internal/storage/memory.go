package storage

import (
	"fmt"
	"strconv"
	"sync"
)

const (
	counterMetric string = "counter"
	gaugeMetric   string = "gauge"
)

type counter struct {
	mx sync.Mutex
	c  map[string]int64
}

type gauge struct {
	mx sync.Mutex
	g  map[string]float64
}

func (c *counter) Load(key string) (int64, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	val, ok := c.c[key]
	return val, ok
}

func (c *counter) LoadAll() map[string]int64 {
	c.mx.Lock()
	defer c.mx.Unlock()
	cCopy := make(map[string]int64)
	for k, v := range c.c {
		cCopy[k] = v
	}
	return cCopy
}

func (c *counter) Store(key string, value int64) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.c[key] += value
}

func (g *gauge) Load(key string) (float64, bool) {
	g.mx.Lock()
	defer g.mx.Unlock()
	val, ok := g.g[key]
	return val, ok
}

func (g *gauge) LoadAll() map[string]float64 {
	g.mx.Lock()
	defer g.mx.Unlock()
	gCopy := make(map[string]float64)
	for k, v := range g.g {
		gCopy[k] = v
	}
	return gCopy
}

func (g *gauge) Store(key string, value float64) {
	g.mx.Lock()
	defer g.mx.Unlock()
	g.g[key] = value
}

type MemStorage struct {
	gauge   gauge
	counter counter
}

func NewMemStorage() MemStorage {
	return MemStorage{gauge: gauge{g: make(map[string]float64)}, counter: counter{c: make(map[string]int64)}}
}

func (storage *MemStorage) Store(metricType, name, value string) error {
	switch metricType {
	case counterMetric:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to convert value %s to an integer: %w", value, err)
		}
		storage.counter.Store(name, v)
		return nil
	case gaugeMetric:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to convert value %s to a float: %w", value, err)
		}
		storage.gauge.Store(name, v)
		return nil
	}
	return fmt.Errorf("metric type does not exist, given type: %v", metricType)
}

func (storage *MemStorage) GetValue(metricType, name string) (string, bool, error) {
	switch metricType {
	case counterMetric:
		v, ok := storage.counter.Load(name)
		if !ok {
			return "", false, nil
		}
		return strconv.FormatInt(v, 10), true, nil
	case gaugeMetric:
		v, ok := storage.gauge.Load(name)
		if !ok {
			return "", false, nil
		}
		return strconv.FormatFloat(v, byte(102), -3, 64), true, nil
	}
	return "", false, fmt.Errorf("metric type does not exist, given type: %v", metricType)
}

func (storage *MemStorage) GetCounterMetrics() map[string]int64 {
	return storage.counter.LoadAll()
}

func (storage *MemStorage) GetGaugeMetrics() map[string]float64 {
	return storage.gauge.LoadAll()
}
