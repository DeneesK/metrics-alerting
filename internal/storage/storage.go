package storage

import (
	"github.com/DeneesK/metrics-alerting/internal/models"
	"go.uber.org/zap"
)

type Storage interface {
	Store(typeMetric string, name string, value interface{}) error
	StoreBanch(metrics []models.Metrics) error
	GetValue(typeMetric, name string) (Result, bool, error)
	GetCounterMetrics() map[string]int64
	GetGaugeMetrics() map[string]float64
	Ping() (bool, error)
}

func NewStorage(filePath string, storeInterval int, isRestore bool, log *zap.SugaredLogger, postgresDSN string) (Storage, error) {
	if postgresDSN != "" {
		return NewDBStorage(postgresDSN, log)
	}

	if filePath != "" {
		return NewFileStorage(filePath, storeInterval, isRestore, log)
	}
	return NewMemStorage(log)
}
