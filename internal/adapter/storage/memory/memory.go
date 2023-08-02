package memory

import "github.com/korovindenis/go-pc-metrics/internal/domain/entity"

type storage struct {
	entity.Metrics
}

func New() (entity.IStorage, error) {
	strg := &storage{}
	strg.Metrics.Gauge = make(map[string]float64)
	strg.Metrics.Counter = make(map[string]int64)

	return strg, nil
}

func (m *storage) SaveGauge(gaugeName string, gaugeValue float64) error {
	m.Metrics.Gauge[gaugeName] = gaugeValue

	return nil
}

func (m *storage) SaveCounter(counterName string, counterValue int64) error {
	m.Metrics.Counter[counterName] = counterValue

	return nil
}

func (m *storage) GetCounter(counterName string) (int64, error) {
	return m.Metrics.Counter[counterName], nil

}
