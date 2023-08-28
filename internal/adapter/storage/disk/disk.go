package disk

import (
	"encoding/json"
	"os"

	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
)

type Storage struct {
	filePath string
}

type cfg interface {
	GetFileStoragePath() string
}

func New(config cfg) (*Storage, error) {
	filePath := config.GetFileStoragePath()

	storage := &Storage{
		filePath: filePath,
	}

	// Создание файла, если он не существует
	//if _, err := os.Stat(filePath); os.IsNotExist(err) {
	if _, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return nil, err
	}
	if err := storage.saveToFile(entity.MetricsType{}); err != nil {
		return nil, err
	}
	//}

	return storage, nil
}

func (m *Storage) SaveGauge(gaugeName string, gaugeValue float64) error {
	metrics, err := m.loadFromFile()
	if err != nil {
		return err
	}

	metrics.Gauge[gaugeName] = gaugeValue
	return m.saveToFile(metrics)
}

func (m *Storage) GetGauge(gaugeName string) (float64, error) {
	metrics, err := m.loadFromFile()
	if err != nil {
		return 0, err
	}

	val, ok := metrics.Gauge[gaugeName]
	if !ok {
		return val, entity.ErrMetricNotFound
	}
	return val, nil
}

func (m *Storage) SaveCounter(counterName string, counterValue int64) error {
	metrics, err := m.loadFromFile()
	if err != nil {
		return err
	}

	metrics.Counter[counterName] = counterValue
	return m.saveToFile(metrics)
}

func (m *Storage) GetCounter(counterName string) (int64, error) {
	metrics, err := m.loadFromFile()
	if err != nil {
		return 0, err
	}

	val, ok := metrics.Counter[counterName]
	if !ok {
		return val, entity.ErrMetricNotFound
	}
	return val, nil
}

func (m *Storage) GetAllData() (entity.MetricsType, error) {
	return m.loadFromFile()
}

func (m *Storage) loadFromFile() (entity.MetricsType, error) {
	file, err := os.OpenFile(m.filePath, os.O_RDONLY, 0666)
	if err != nil {
		return entity.MetricsType{}, err
	}
	defer file.Close()

	var metrics entity.MetricsType
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&metrics)
	if err != nil {
		return entity.MetricsType{}, err
	}

	return metrics, nil
}

func (m *Storage) saveToFile(metrics entity.MetricsType) error {
	file, err := os.OpenFile(m.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(metrics)
	if err != nil {
		return err
	}

	return nil
}
