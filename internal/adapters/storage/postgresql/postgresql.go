package postgresql

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	"github.com/pressly/goose"
)

type Storage struct {
	db *sql.DB
}

type cfg interface {
	GetFileStoragePath() string
	GetRestore() bool
	GetDatabaseConnectionString() string
}

func New(config cfg) (*Storage, error) {
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
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	for _, v := range metrics {
		switch v.MType {
		case "gauge":
			_, err = tx.ExecContext(ctx, `INSERT INTO gauge (name, value) VALUES ($1, $2)`, v.ID, v.Value)
		case "counter":
			_, err = tx.ExecContext(ctx, `INSERT INTO counter (name, delta) VALUES ($1, $2)`, v.ID, v.Delta)
		default:
			tx.Rollback()
			return entity.ErrInputVarIsWrongType
		}
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *Storage) SaveGauge(ctx context.Context, gaugeName string, gaugeValue float64) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO gauge (name, value) VALUES ($1, $2)`, gaugeName, gaugeValue)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetGauge(ctx context.Context, gaugeName string) (float64, error) {
	var gaugeValue float64
	err := s.db.QueryRowContext(ctx, "SELECT value FROM gauge WHERE name = $1 ORDER BY created DESC", gaugeName).Scan(&gaugeValue)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, entity.ErrMetricNotFound
		} else {
			return 0, err
		}
	}
	return gaugeValue, nil
}

func (s *Storage) SaveCounter(ctx context.Context, counterName string, counterValue int64) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO counter (name, delta) VALUES ($1, $2);`, counterName, counterValue)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetCounter(ctx context.Context, counterName string) (int64, error) {
	var counterValue int64
	err := s.db.QueryRowContext(ctx, "SELECT delta FROM counter WHERE name = $1 ORDER BY created DESC", counterName).Scan(&counterValue)
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
		ORDER BY name, created DESC;
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
		ORDER BY name, created DESC;
	`)
	if err != nil {
		return metrics, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var value int64
		if err := rows.Scan(&name, &value); err != nil {
			return metrics, err
		}
		if err := rows.Err(); err != nil {
			return metrics, err
		}
		metrics.Counter[name] = value
	}

	return metrics, nil
}

func (s *Storage) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return s.db.PingContext(ctx)
}
