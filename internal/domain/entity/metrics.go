package entity

type (
	GaugeType   map[string]float64
	CounterType map[string]int64
)

// types of metrics
type MetricsType struct {
	Gauge   GaugeType
	Counter CounterType
}

// for get requests in the url
type MetricsURI struct {
	MetricType string
	MetricName string
	MetricVal  string
}

// app metrics
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
