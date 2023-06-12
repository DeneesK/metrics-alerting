package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/DeneesK/metrics-alerting/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

const createTableQuery = `CREATE TABLE IF NOT EXISTS metrics(
	"metrictype" TEXT NOT NULL,
	"metricname" TEXT NOT NULL UNIQUE,
	"counter" INTEGER,
	"gauge" DOUBLE PRECISION)`

type DBStorage struct {
	log *zap.SugaredLogger
	db  *sql.DB
}

func NewDBStorage(postgresDSN string, log *zap.SugaredLogger) (*DBStorage, error) {
	db, err := NewDBSession(postgresDSN)
	if err != nil {
		return nil, fmt.Errorf("during initializing of new db session, error occurred: %w", err)
	}

	err = createTable(db)
	if err != nil {
		return nil, fmt.Errorf("impossible to create table: %w", err)
	}

	return &DBStorage{log: log, db: db}, nil
}

func NewDBSession(postgresDSN string) (*sql.DB, error) {
	db, err := sql.Open("pgx", postgresDSN)
	if err != nil {
		for i, atmp := range readAttempts {
			time.Sleep(atmp)
			db, err = sql.Open("pgx", postgresDSN)
			if err != nil && i < 2 {
				continue
			}
			if err != nil && i == 2 {
				return nil, fmt.Errorf("unable to connect to db: %w", err)
			}
		}
	}
	return db, nil
}

func (storage *DBStorage) Ping() (bool, error) {
	err := storage.db.Ping()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (storage *DBStorage) Store(typeMetric string, name string, value interface{}) error {
	switch typeMetric {
	case counterMetric:
		v, ok := value.(int64)
		if !ok {
			return fmt.Errorf("value cannot be cast to a specific type")
		}
		_, err := storage.db.Exec("INSERT INTO metrics VALUES ($1, $2, $3, $4) ON CONFLICT (metricname) DO UPDATE SET counter=metrics.counter+$3", counterMetric, name, v, 0)
		if err != nil {
			return fmt.Errorf("during attempt to store data to database error ocurred: %w", err)
		}
	case gaugeMetric:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("value cannot be cast to a specific type")
		}
		_, err := storage.db.Exec("INSERT INTO metrics VALUES ($1, $2, $3, $4) ON CONFLICT (metricname) DO UPDATE SET gauge=$4", gaugeMetric, name, 0, v)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("metric type does not exist, given type: %v", typeMetric)
	}
	return nil
}

func (storage *DBStorage) StoreBanch(metrics []models.Metrics) error {
	ctx := context.Background()
	if len(metrics) < 1001 {
		if err := storage.insertBanch(ctx, metrics); err != nil {
			return fmt.Errorf("postgres db error: %w", err)
		}
		return nil
	}
	banch := make([]models.Metrics, 0, 1000)
	for _, m := range metrics {
		banch = append(banch, m)
		if len(banch) == 1000 {
			if err := storage.insertBanch(ctx, metrics); err != nil {
				return fmt.Errorf("postgres db error: %w", err)
			}
			banch = banch[:0]
		}
	}
	if err := storage.insertBanch(ctx, banch); err != nil {
		return fmt.Errorf("postgres db error: %w", err)
	}
	return nil
}

func (storage *DBStorage) GetValue(typeMetric, name string) (Result, bool, error) {
	switch typeMetric {
	case counterMetric:
		var value int64
		row := storage.db.QueryRowContext(context.Background(), "SELECT metrics.counter FROM metrics WHERE metrics.MetricName=$1", name)
		err := row.Scan(&value)
		if err != nil {
			return Result{}, false, err
		}
		return Result{Counter: value, Gauge: 0}, true, nil
	case gaugeMetric:
		var value float64
		row := storage.db.QueryRowContext(context.Background(), "SELECT metrics.gauge FROM metrics WHERE metrics.MetricName=$1", name)
		err := row.Scan(&value)
		if err != nil {
			return Result{}, false, err
		}
		return Result{Counter: 0, Gauge: value}, true, nil
	default:
		return Result{}, false, fmt.Errorf("metric type does not exist, given type: %v", typeMetric)
	}
}

func (storage *DBStorage) GetCounterMetrics() map[string]int64 {
	return make(map[string]int64, 0)
}

func (storage *DBStorage) GetGaugeMetrics() map[string]float64 {
	return make(map[string]float64, 0)
}

func (storage *DBStorage) insertBanch(ctx context.Context, metrics []models.Metrics) error {
	tx, err := storage.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO metrics VALUES ($1, $2, $3, $4) ON CONFLICT (metricname) DO UPDATE SET gauge=$4, counter=metrics.counter+$3")
	if err != nil {
		return fmt.Errorf("insert error: %w", err)
	}
	defer stmt.Close()
	for _, m := range metrics {
		_, err := stmt.ExecContext(ctx, m.MType, m.ID, m.Delta, m.Value)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func createTable(session *sql.DB) error {
	_, err := session.Exec(createTableQuery)
	if err != nil {
		return err
	}
	return nil
}
