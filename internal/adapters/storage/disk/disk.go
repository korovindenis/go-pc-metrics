// Storage with file system
package disk

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	"go.uber.org/zap/zapcore"
)

type Storage struct {
	filePath string
	metrics  entity.MetricsType
}

type cfg interface {
	GetFileStoragePath() string
	GetRestore() bool
}

type log interface {
	Info(msg string, fields ...zapcore.Field)
}

func New(config cfg, log log) (*Storage, error) {
	log.Info("Storage is disk")

	storage := &Storage{
		filePath: config.GetFileStoragePath(),
		metrics: entity.MetricsType{
			Gauge:   make(entity.GaugeType),
			Counter: make(entity.CounterType),
		},
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

func (s *Storage) SaveAllData(ctx context.Context, metrics []entity.Metrics) error {
	return s.saveToFile()
}

func (s *Storage) SaveGauge(ctx context.Context, gaugeName string, gaugeValue float64) error {
	s.metrics.Gauge[gaugeName] = gaugeValue
	return nil
}

func (s *Storage) GetGauge(ctx context.Context, gaugeName string) (float64, error) {
	val, ok := s.metrics.Gauge[gaugeName]
	if !ok {
		return val, entity.ErrMetricNotFound
	}
	return val, nil
}

func (s *Storage) SaveCounter(ctx context.Context, counterName string, counterValue int64) error {
	s.metrics.Counter[counterName] = counterValue
	return nil
}

func (s *Storage) GetCounter(ctx context.Context, counterName string) (int64, error) {
	val, ok := s.metrics.Counter[counterName]
	if !ok {
		return val, entity.ErrMetricNotFound
	}
	return val, nil
}

func (s *Storage) GetAllData(ctx context.Context) (entity.MetricsType, error) {
	return s.metrics, nil
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

func (s *Storage) Ping(ctx context.Context) error {
	return nil
}
