package memory

import (
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
)

type Storage struct {
	entity.MetricsType
}

func New() (*Storage, error) {
	strg := &Storage{}
	strg.MetricsType.Gauge = make(map[string]float64)
	strg.MetricsType.Counter = make(map[string]int64)

	return strg, nil
}

func (m *Storage) SaveGauge(gaugeName string, gaugeValue float64) error {
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
	return entity.MetricsType{
		Gauge:   m.MetricsType.Gauge,
		Counter: m.MetricsType.Counter,
	}, nil
}
