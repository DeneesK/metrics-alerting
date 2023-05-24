package storage

import (
	"fmt"
	"os"
	"sync"
)

const (
	counterMetric string = "counter"
	gaugeMetric   string = "gauge"
)

type Result struct {
	Counter int64
	Gauge   float64
}

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

func (c *counter) Set(m map[string]int64) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.c = m
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

func (g *gauge) Set(m map[string]float64) {
	g.mx.Lock()
	defer g.mx.Unlock()
	g.g = m
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

func (storage *MemStorage) Store(metricType, name string, value interface{}) error {
	switch metricType {
	case counterMetric:
		v := value.(int64)
		storage.counter.Store(name, v)
		return nil
	case gaugeMetric:
		v := value.(float64)
		storage.gauge.Store(name, v)
		return nil
	}
	return fmt.Errorf("metric type does not exist, given type: %v", metricType)
}

func (storage *MemStorage) GetValue(metricType, name string) (Result, bool, error) {
	switch metricType {
	case counterMetric:
		v, ok := storage.counter.Load(name)
		if !ok {
			return Result{0, 0}, false, nil
		}
		return Result{Counter: v, Gauge: 0}, true, nil
	case gaugeMetric:
		v, ok := storage.gauge.Load(name)
		if !ok {
			return Result{0, 0}, false, nil
		}
		return Result{Counter: 0, Gauge: v}, true, nil
	}
	return Result{0, 0}, false, fmt.Errorf("metric type does not exist, given type: %v", metricType)
}

func (storage *MemStorage) GetCounterMetrics() map[string]int64 {
	return storage.counter.LoadAll()
}

func (storage *MemStorage) GetGaugeMetrics() map[string]float64 {
	return storage.gauge.LoadAll()
}

func (storage *MemStorage) setMetrics(metrics *allMetrics) {
	storage.counter.Set(metrics.Counter)
	storage.gauge.Set(metrics.Gauge)
}

func (storage *MemStorage) SaveToFile(path string) error {
	return storeToFile(path, storage)
}

func (storage *MemStorage) LoadFromFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	return loadFromFile(path, storage)
}
