package memory

import (
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	serverusecase "github.com/korovindenis/go-pc-metrics/internal/domain/usecase/server"
)

type storage struct {
	entity.MetricsType
}

func New() (serverusecase.IStorage, error) {
	strg := &storage{}
	strg.MetricsType.Gauge = make(map[string]float64)
	strg.MetricsType.Counter = make(map[string]int64)

	return strg, nil
}

func (m *storage) SaveGauge(gaugeName string, gaugeValue float64) error {
	m.MetricsType.Gauge[gaugeName] = gaugeValue

	return nil
}

func (m *storage) SaveCounter(counterName string, counterValue int64) error {
	m.MetricsType.Counter[counterName] = counterValue

	return nil
}

func (m *storage) GetCounter(counterName string) (int64, error) {
	return m.MetricsType.Counter[counterName], nil

}
