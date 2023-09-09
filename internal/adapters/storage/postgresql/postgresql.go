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

func (s *Storage) SaveAllData() error {
	return nil
}

func (s *Storage) SaveGauge(gaugeName string, gaugeValue float64) error {
	_, err := s.db.Exec(`
		INSERT INTO gauge (name, value)
		VALUES ($1, $2)
		ON CONFLICT (name)
		DO UPDATE SET value = EXCLUDED.value;
	`, gaugeName, gaugeValue)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetGauge(gaugeName string) (float64, error) {
	var gaugeValue float64
	err := s.db.QueryRow("SELECT value FROM gauge WHERE name = $1", gaugeName).Scan(&gaugeValue)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, entity.ErrMetricNotFound
		} else {
			return 0, err
		}
	}
	return gaugeValue, nil
}

func (s *Storage) SaveCounter(counterName string, counterValue int64) error {
	_, err := s.db.Exec(`
		INSERT INTO counter (name, value)
		VALUES ($1, $2)
		ON CONFLICT (name)
		DO UPDATE SET value = EXCLUDED.value;
	`, counterName, counterValue)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetCounter(counterName string) (int64, error) {
	var counterValue int64
	err := s.db.QueryRow("SELECT value FROM counter WHERE name = $1", counterName).Scan(&counterValue)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, entity.ErrMetricNotFound
		} else {
			return 0, err
		}
	}
	return counterValue, nil
}

func (s *Storage) GetAllData() (entity.MetricsType, error) {
	metrics := entity.MetricsType{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}

	rows, err := s.db.Query("SELECT name, value FROM gauge")
	if err != nil {
		return metrics, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var value float64
		err := rows.Scan(&name, &value)
		if err != nil {
			return metrics, err
		}
		if err := rows.Err(); err != nil {
			return metrics, err
		}
		metrics.Gauge[name] = value
	}

	rows, err = s.db.Query("SELECT name, value FROM counter")
	if err != nil {
		return metrics, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var value int64
		err := rows.Scan(&name, &value)
		if err != nil {
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
