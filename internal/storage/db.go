package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/DeneesK/metrics-alerting/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sethvargo/go-retry"
	"go.uber.org/zap"
)

const createTableQuery = `CREATE TABLE IF NOT EXISTS metrics(
	"metrictype" TEXT NOT NULL,
	"metricname" TEXT NOT NULL UNIQUE,
	"counter" BIGINT,
	"gauge" DOUBLE PRECISION)`

type DBStorage struct {
	log *zap.SugaredLogger
	db  *sql.DB
}

func NewDBStorage(dsn string, log *zap.SugaredLogger) (*DBStorage, error) {
	db, err := NewDBSession(dsn)
	if err != nil {
		return nil, fmt.Errorf("during initializing of new db session, error occurred: %w", err)
	}

	err = createTable(db)
	if err != nil {
		return nil, fmt.Errorf("impossible to create table: %w", err)
	}

	return &DBStorage{log: log, db: db}, nil
}

func NewDBSession(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	b := retry.WithMaxRetries(3, retry.NewExponential(1*time.Second))
	err = retry.Do(ctx, b, try(db))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (storage *DBStorage) Ping() error {
	err := storage.db.Ping()
	if err != nil {
		return err
	}
	return nil
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

func (storage *DBStorage) StoreBatch(metrics []models.Metrics) error {
	ctx := context.Background()
	start := 0
	end := 1000
	for len(metrics) > start {
		if len(metrics)-end < 0 {
			if err := storage.insertBatch(ctx, metrics[start:]); err != nil {
				return fmt.Errorf("postgres db error: %w", err)
			}
			break
		}
		if err := storage.insertBatch(ctx, metrics[start:end]); err != nil {
			return fmt.Errorf("postgres db error: %w", err)
		}
		start += 1000
		end += 1000
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

func (storage *DBStorage) GetCounterMetrics() (map[string]int64, error) {
	metrics := make(map[string]int64, 0)
	var name string
	var v int64
	rows, err := storage.db.QueryContext(context.Background(), "SELECT metrics.metricname, metrics.counter FROM metrics WHERE metrics.metrictype=$1", counterMetric)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name, &v)
		if err != nil {
			return nil, err
		}
		metrics[name] = v
	}
	return metrics, nil
}

func (storage *DBStorage) GetGaugeMetrics() (map[string]float64, error) {
	metrics := make(map[string]float64, 0)
	var name string
	var v float64
	rows, err := storage.db.QueryContext(context.Background(), "SELECT metrics.metricname, metrics.gauge FROM metrics WHERE metrics.metrictype=$1", gaugeMetric)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name, &v)
		if err != nil {
			return nil, err
		}
		metrics[name] = v
	}
	return metrics, nil
}

func (storage *DBStorage) Close() error {
	return storage.db.Close()
}

func (storage *DBStorage) insertBatch(ctx context.Context, metrics []models.Metrics) error {
	tx, err := storage.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	values := make(map[string][]interface{}, len(metrics))
	var valString []string
	var v []interface{}
	for _, m := range metrics {
		values[m.ID] = []interface{}{m.MType, m.ID, m.Delta, m.Value}
	}
	i := 0
	for _, m := range values {
		valString = append(valString, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		v = append(v, m[0])
		v = append(v, m[1])
		v = append(v, m[2])
		v = append(v, m[3])
		i++
	}
	smt := "INSERT INTO metrics VALUES %s ON CONFLICT (metricname) DO UPDATE SET gauge=EXCLUDED.gauge, counter=metrics.counter+EXCLUDED.counter"
	smt = fmt.Sprintf(smt, strings.Join(valString, ","))
	_, err = tx.ExecContext(ctx, smt, v...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("insert batch error: %w", err)
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

func try(db *sql.DB) func(context.Context) error {
	return func(ctx context.Context) error {
		if err := db.PingContext(ctx); err != nil {
			return retry.RetryableError(err)
		}
		return nil
	}
}
