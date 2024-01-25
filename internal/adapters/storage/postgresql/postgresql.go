// Storage with postgresql
package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	"github.com/lib/pq"
	"github.com/pressly/goose"
	"go.uber.org/zap/zapcore"
)

const (
	UniqueViolation = "unique_violation"
)

type Storage struct {
	db *sql.DB
}

//go:generate mockery --name log --exported
type log interface {
	Info(msg string, fields ...zapcore.Field)
}

//go:generate mockery --name cfg --exported
type cfg interface {
	GetFileStoragePath() string
	GetRestore() bool
	GetDatabaseConnectionString() string
}

func New(config cfg, log log) (*Storage, error) {
	log.Info("Storage is database")

	db, err := sql.Open("pgx", config.GetDatabaseConnectionString())
	if err != nil {
		return nil, err
	}

	storage := &Storage{
		db: db,
	}

	if err := storage.runMigrations(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) runMigrations() error {
	if err := goose.Run("up", s.db, "deployments/db/migrations"); err != nil {
		return err
	}

	return nil
}

func (s *Storage) SaveAllData(ctx context.Context, metrics []entity.Metrics) error {
	var err error

	for _, v := range metrics {
		switch v.MType {
		case "gauge":
			query := "INSERT INTO gauge (name, value) VALUES ($1,$2)"
			err = s.retryableExec(ctx, query, v.ID, v.Value)
		case "counter":
			query := "INSERT INTO counter (name, delta) VALUES ($1,$2)"
			err = s.retryableExec(ctx, query, v.ID, v.Delta)
		default:
			return entity.ErrInputVarIsWrongType
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) SaveGauge(ctx context.Context, gaugeName string, gaugeValue float64) error {
	query := "INSERT INTO gauge (name, value) VALUES ($1, $2)"
	if err := s.retryableExec(ctx, query, gaugeName, gaugeValue); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetGauge(ctx context.Context, gaugeName string) (float64, error) {
	var gaugeValue float64
	err := s.db.QueryRowContext(ctx, "SELECT value FROM gauge WHERE name = $1 ORDER BY id DESC", gaugeName).Scan(&gaugeValue)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, entity.ErrMetricNotFound
		}
		return 0, err
	}
	return gaugeValue, nil
}

func (s *Storage) SaveCounter(ctx context.Context, counterName string, counterValue int64) error {
	query := "INSERT INTO counter (name, delta) VALUES ($1, $2)"
	if err := s.retryableExec(ctx, query, counterName, counterValue); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetCounter(ctx context.Context, counterName string) (int64, error) {
	var counterValue int64
	err := s.db.QueryRowContext(ctx, "SELECT delta FROM counter WHERE name = $1 ORDER BY id DESC", counterName).Scan(&counterValue)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, entity.ErrMetricNotFound
		} else {
			return 0, err
		}
	}
	return counterValue, nil
}

func (s *Storage) GetAllData(ctx context.Context) (entity.MetricsType, error) {
	metrics := entity.MetricsType{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT DISTINCT ON (name) name, value
		FROM gauge
		ORDER BY name, id DESC;
	`)
	if err != nil {
		return metrics, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var value float64
		if err := rows.Scan(&name, &value); err != nil {
			return metrics, err
		}
		if err := rows.Err(); err != nil {
			return metrics, err
		}
		metrics.Gauge[name] = value
	}

	rows, err = s.db.QueryContext(ctx, `
		SELECT DISTINCT ON (name) name, delta
		FROM counter
		ORDER BY name, id DESC;
	`)
	if err != nil {
		return metrics, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var delta int64
		if err := rows.Scan(&name, &delta); err != nil {
			return metrics, err
		}
		if err := rows.Err(); err != nil {
			return metrics, err
		}
		metrics.Counter[name] = delta
	}

	return metrics, nil
}

func (s *Storage) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return s.db.PingContext(ctx)
}

func (s *Storage) retryableExec(ctx context.Context, query string, args ...any) error {
	const maxRetries = 3
	var retryDelays = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	var pqError *pq.Error

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 0; i <= maxRetries; i++ {
		if i > 0 {
			// Waiting before trying again
			time.Sleep(retryDelays[i-1])
		}

		if _, err = stmt.ExecContext(ctx, args...); err != nil {
			// Checking for a unique violation (UniqueViolation)
			if errors.As(err, &pqError) {
				if pqError.Code.Name() == UniqueViolation {
					continue
				}
			}

			return err
		}
		return tx.Commit() // Successful execution
	}
	return errors.New("max retries exceeded")
}
