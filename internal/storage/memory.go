package storage

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

const (
	counterMetric string = "counter"
	gaugeMetric   string = "gauge"
)

var ErrMetricType = errors.New("metric type does not exist")
var ErrCounterValue = errors.New("wrong counter value, must be integer")
var ErrGaugeValue = errors.New("wrong gauge value, must be float")

type counter struct {
	mx sync.Mutex
	c  map[string]int
}

type gauge struct {
	mx sync.Mutex
	g  map[string]float64
}

func (c *counter) Load(key string) (int, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	val, ok := c.c[key]
	return val, ok
}

func (c *counter) LoadAll() map[string]int {
	c.mx.Lock()
	defer c.mx.Unlock()
	return c.c
}

func (c *counter) Store(key string, value int) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.c[key] = value
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
	return g.g
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
	return MemStorage{gauge: gauge{g: make(map[string]float64)}, counter: counter{c: make(map[string]int)}}
}

func (storage *MemStorage) Store(type_, name, value string) error {
	switch type_ {
	case counterMetric:
		v, err := strconv.Atoi(value)
		if err != nil {
			return ErrCounterValue
		}
		storage.counter.Store(name, v)
		return nil
	case gaugeMetric:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return ErrGaugeValue
		}

		storage.gauge.Store(name, v)
		return nil
	}
	return ErrMetricType
}

func (storage *MemStorage) GetValue(typeMetric, name string) (string, bool, error) {
	switch typeMetric {
	case counterMetric:
		v, ok := storage.counter.Load(name)
		if !ok {
			return "", false, nil
		}
		return fmt.Sprintf("%d", v), true, nil
	case gaugeMetric:
		v, ok := storage.gauge.Load(name)
		if !ok {
			return "", false, nil
		}
		return fmt.Sprintf("%g", v), true, nil
	}
	return "", false, ErrMetricType
}

func (storage *MemStorage) GetCounterMetrics() map[string]int {
	return storage.counter.LoadAll()
}

func (storage *MemStorage) GetGaugeMetrics() map[string]float64 {
	return storage.gauge.LoadAll()
}
