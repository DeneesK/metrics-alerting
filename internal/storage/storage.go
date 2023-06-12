package storage

import (
	"time"

	"github.com/DeneesK/metrics-alerting/internal/models"
	"go.uber.org/zap"
)

const (
	fstAttempt   time.Duration = time.Duration(0) * time.Second
	sndAttempt   time.Duration = fstAttempt * 1
	thirdAttempt time.Duration = fstAttempt * 1
)

var readAttempts = []time.Duration{fstAttempt, sndAttempt, thirdAttempt}

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
