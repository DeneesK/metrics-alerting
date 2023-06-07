package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

const createTablestring = `CREATE TABLE IF NOT EXISTS metrics(
	"metrictype" TEXT,
	"metricname" TEXT,
	"counter" INTEGER,
	"gauge" DOUBLE PRECISION)`

type DBStorage struct {
	log *zap.SugaredLogger
	db  *sql.DB
}

func NewDBStorage(postgresDSN string, log *zap.SugaredLogger) (*DBStorage, error) {
	db, err := NewDBSession(postgresDSN)
	if err != nil {
		return nil, fmt.Errorf("during initializing of new db session, error occurred: %v", err)
	}

	err = createTable(db)
	if err != nil {
		return nil, fmt.Errorf("impossible to create table: %v", err)
	}

	return &DBStorage{log: log, db: db}, nil
}

func NewDBSession(postgresDSN string) (*sql.DB, error) {
	db, err := sql.Open("pgx", postgresDSN)
	if err != nil {
		return nil, err
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
		_, err := storage.db.Exec("INSERT INTO metrics VALUES ($1, $2, $3, $4)", counterMetric, name, v, 0)
		if err != nil {
			return err
		}
	case gaugeMetric:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("value cannot be cast to a specific type")
		}
		_, err := storage.db.Exec("INSERT INTO metrics VALUES ($1, $2, $3, $4)", gaugeMetric, name, 0, v)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("metric type does not exist, given type: %v", typeMetric)
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

func createTable(session *sql.DB) error {
	_, err := session.Exec(createTablestring)
	if err != nil {
		return err
	}
	return nil
}
