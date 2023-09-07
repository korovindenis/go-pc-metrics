package postgresql

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
)

type Storage struct {
	filePath string
	metrics  entity.MetricsType
	db       *sql.DB
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
		filePath: config.GetFileStoragePath(),
		metrics: entity.MetricsType{
			Gauge:   make(entity.GaugeType),
			Counter: make(entity.CounterType),
		},
		db: db,
	}
	// create dir
	if err := os.MkdirAll(filepath.Dir(storage.filePath), 0770); err != nil {
		return nil, err
	}
	// open file
	file, err := os.OpenFile(storage.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if config.GetRestore() {
		metrics, err := storage.loadFromFile()
		if err != nil {
			return storage, nil
		}
		storage.metrics = metrics
	}

	return storage, nil
}

func (s *Storage) SaveAllData() error {
	return s.saveToFile()
}

func (s *Storage) SaveGauge(gaugeName string, gaugeValue float64) error {
	s.metrics.Gauge[gaugeName] = gaugeValue
	return nil
}

func (s *Storage) GetGauge(gaugeName string) (float64, error) {
	val, ok := s.metrics.Gauge[gaugeName]
	if !ok {
		return val, entity.ErrMetricNotFound
	}
	return val, nil
}

func (s *Storage) SaveCounter(counterName string, counterValue int64) error {
	s.metrics.Counter[counterName] = counterValue
	return nil
}

func (s *Storage) GetCounter(counterName string) (int64, error) {
	val, ok := s.metrics.Counter[counterName]
	if !ok {
		return val, entity.ErrMetricNotFound
	}
	return val, nil
}

func (s *Storage) GetAllData() (entity.MetricsType, error) {
	return s.metrics, nil
}

func (s *Storage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.db.PingContext(ctx)
}

func (s *Storage) loadFromFile() (entity.MetricsType, error) {
	file, err := os.OpenFile(s.filePath, os.O_RDONLY, 0666)
	if err != nil {
		return entity.MetricsType{
			Gauge:   make(entity.GaugeType),
			Counter: make(entity.CounterType),
		}, err
	}
	defer file.Close()

	var metrics entity.MetricsType
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&metrics)
	if err != nil {
		return entity.MetricsType{
			Gauge:   make(entity.GaugeType),
			Counter: make(entity.CounterType),
		}, err
	}

	return metrics, nil
}

func (s *Storage) saveToFile() error {
	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(s.metrics)
	if err != nil {
		return err
	}

	return nil
}
