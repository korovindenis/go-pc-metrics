package memory

import (
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
)

type Storage struct {
	MetricsType entity.MetricsType
}

type cfg interface {
}

func New(config cfg) (*Storage, error) {
	storage := Storage{
		MetricsType: entity.MetricsType{
			Gauge:   make(map[string]float64),
			Counter: make(map[string]int64),
		},
	}

	return &storage, nil
}

func (m *Storage) SaveGauge(gaugeName string, gaugeValue float64) error {
	//if m.MetricsType.Gauge == nil {
	//m.MetricsType.Gauge = make(map[string]float64)
	//}

	m.MetricsType.Gauge[gaugeName] = gaugeValue

	return nil
}

func (m *Storage) GetGauge(gaugeName string) (float64, error) {
	val, ok := m.MetricsType.Gauge[gaugeName]
	if !ok {
		return val, entity.ErrMetricNotFound
	}
	return val, nil
}

func (m *Storage) SaveCounter(counterName string, counterValue int64) error {
	//if m.MetricsType.Counter == nil {
	//m.MetricsType.Counter = make(map[string]int64)
	//}

	m.MetricsType.Counter[counterName] = counterValue

	return nil
}

func (m *Storage) GetCounter(counterName string) (int64, error) {
	val, ok := m.MetricsType.Counter[counterName]
	if !ok {
		return val, entity.ErrMetricNotFound
	}
	return val, nil

}

func (m *Storage) GetAllData() (entity.MetricsType, error) {
	return m.MetricsType, nil
}

func (m *Storage) SaveAllData() error {
	return nil
}
