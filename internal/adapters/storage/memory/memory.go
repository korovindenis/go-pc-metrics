package memory

import (
	"context"

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

func (m *Storage) SaveGauge(ctx context.Context, gaugeName string, gaugeValue float64) error {
	m.MetricsType.Gauge[gaugeName] = gaugeValue

	return nil
}

func (m *Storage) GetGauge(ctx context.Context, gaugeName string) (float64, error) {
	val, ok := m.MetricsType.Gauge[gaugeName]
	if !ok {
		return val, entity.ErrMetricNotFound
	}
	return val, nil
}

func (m *Storage) SaveCounter(ctx context.Context, counterName string, counterValue int64) error {
	m.MetricsType.Counter[counterName] = counterValue

	return nil
}

func (m *Storage) GetCounter(ctx context.Context, counterName string) (int64, error) {
	val, ok := m.MetricsType.Counter[counterName]
	if !ok {
		return val, entity.ErrMetricNotFound
	}
	return val, nil

}

func (m *Storage) GetAllData(ctx context.Context) (entity.MetricsType, error) {
	return m.MetricsType, nil
}

func (m *Storage) SaveAllData(ctx context.Context, metrics []entity.Metrics) error {
	return nil
}

func (m *Storage) Ping(ctx context.Context) error {
	return nil
}
