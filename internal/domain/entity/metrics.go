package entity

type GaugeType map[string]float64
type CounterType map[string]int64

type MetricsType struct {
	Gauge   GaugeType
	Counter CounterType
}
