package entity

type (
	GaugeType   map[string]float64
	CounterType map[string]int64
)

type MetricsType struct {
	Gauge   GaugeType
	Counter CounterType
}

// type MetricsURI struct {
// 	MetricType string
// 	MetricName string
// 	MetricVal  string
// }

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
