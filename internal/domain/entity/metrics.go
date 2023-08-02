package entity

type Metrics struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

type IStorage interface {
	SaveGauge(gaugeName string, gaugeValue float64) error

	SaveCounter(counterName string, counterValue int64) error
	GetCounter(counterName string) (int64, error)
}
