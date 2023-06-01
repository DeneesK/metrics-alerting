package storage

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"go.uber.org/zap"
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

func (c *counter) set(m map[string]int64) {
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

func (g *gauge) set(m map[string]float64) {
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
	gauge         gauge
	counter       counter
	filePath      string
	storeInterval time.Duration
	log           *zap.SugaredLogger
}

func NewMemStorage(filePath string, storeInterval int, isRestore bool, log *zap.SugaredLogger) *MemStorage {
	ms := MemStorage{
		gauge:         gauge{g: make(map[string]float64)},
		counter:       counter{c: make(map[string]int64)},
		filePath:      filePath,
		storeInterval: time.Duration(storeInterval) * time.Second,
		log:           log,
	}

	if filePath != "" {
		if isRestore {
			if err := ms.loadFromFile(filePath); err != nil {
				ms.log.Errorf("during attempt to load from file, error occurred: %v", err)
			}
		}
		if storeInterval != 0 {
			if err := ms.startStoring(); err != nil {
				ms.log.Errorf("during initializing of new storage, error occurred: %v", err)
			}
		}
	}

	return &ms
}

func (storage *MemStorage) Store(metricType, name string, value interface{}) error {
	switch metricType {
	case counterMetric:
		v, ok := value.(int64)
		if !ok {
			return fmt.Errorf("value cannot be cast to a specific type")
		}
		storage.counter.Store(name, v)
		return nil
	case gaugeMetric:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("value cannot be cast to a specific type")
		}
		storage.gauge.Store(name, v)
		return nil
	default:
		return fmt.Errorf("metric type does not exist, given type: %v", metricType)
	}
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

func (storage *MemStorage) saveToFile(path string) error {
	return storeToFile(path, storage)
}

func (storage *MemStorage) loadFromFile(path string) error {
	return loadFromFile(path, storage)
}

func (storage *MemStorage) setMetrics(metrics *allMetrics) {
	storage.counter.set(metrics.Counter)
	storage.gauge.set(metrics.Gauge)
}

func (storage *MemStorage) startStoring() error {
	dir, _ := path.Split(storage.filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0666)
		if err != nil {
			return err
		}
	}
	go storage.save()
	return nil
}

func (storage *MemStorage) save() {
	for {
		time.Sleep(storage.storeInterval)
		if err := storage.saveToFile(storage.filePath); err != nil {
			storage.log.Errorf("during attempt to store data to file, error occurred: %v", err)
		}
	}
}
