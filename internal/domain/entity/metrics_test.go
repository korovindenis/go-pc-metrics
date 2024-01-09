package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	tests := []struct {
		name          string
		metric        Metrics
		expectedID    string
		expectedMType string
		expectedDelta *int64
		expectedValue *float64
	}{
		{
			name: "Metrics with Delta",
			metric: Metrics{
				ID:    "123",
				MType: "gauge",
				Delta: new(int64),
				Value: nil,
			},
			expectedID:    "123",
			expectedMType: "gauge",
			expectedDelta: new(int64),
			expectedValue: nil,
		},
		{
			name: "Metrics with Value",
			metric: Metrics{
				ID:    "456",
				MType: "counter",
				Delta: nil,
				Value: new(float64),
			},
			expectedID:    "456",
			expectedMType: "counter",
			expectedDelta: nil,
			expectedValue: new(float64),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedID, tt.metric.ID)
			assert.Equal(t, tt.expectedMType, tt.metric.MType)
			assert.Equal(t, tt.expectedDelta, tt.metric.Delta)
			assert.Equal(t, tt.expectedValue, tt.metric.Value)
		})
	}
}
