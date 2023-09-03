package disk

import (
	"encoding/json"
	"os"

	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
)

type Storage struct {
	isRestore bool
	filePath  string
	metrics   entity.MetricsType
}

type cfg interface {
	GetFileStoragePath() string
	GetRestore() bool
}

func New(config cfg) (*Storage, error) {
	storage := &Storage{
		filePath:  config.GetFileStoragePath(),
		isRestore: config.GetRestore(),
		metrics: entity.MetricsType{
			Gauge:   make(entity.GaugeType),
			Counter: make(entity.CounterType),
		},
	}

	if _, err := os.OpenFile(storage.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return nil, err
	}

	// if storage.isRestore {
	// 	metrics, err := storage.loadFromFile()
	// 	if err == nil {
	// 		storage.metrics = metrics
	// 	}
	// }

	return storage, nil
}

func (s *Storage) SaveAllData() error {
	return s.saveToFile()
}

func (s *Storage) SaveGauge(gaugeName string, gaugeValue float64) error {
	s.metrics.Gauge[gaugeName] = gaugeValue
	return s.saveToFile()
}

func (s *Storage) GetGauge(gaugeName string) (float64, error) {
	val, ok := s.metrics.Gauge[gaugeName]
	if !ok {
		return val, entity.ErrMetricNotFound
	}
	return val, nil
}

func (s *Storage) SaveCounter(counterName string, counterValue int64) error {
	// metrics, err := s.loadFromFile()
	// if err != nil {
	// 	return err
	// }

	s.metrics.Counter[counterName] = counterValue
	return s.saveToFile()
}

func (s *Storage) GetCounter(counterName string) (int64, error) {
	// metrics, err := s.loadFromFile()
	// if err != nil {
	// 	return 0, err
	// }

	val, ok := s.metrics.Counter[counterName]
	if !ok {
		return val, entity.ErrMetricNotFound
	}
	return val, nil
}

func (s *Storage) GetAllData() (entity.MetricsType, error) {
	return s.loadFromFile()
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
