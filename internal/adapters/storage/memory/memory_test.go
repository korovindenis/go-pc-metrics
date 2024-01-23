// Storage with RAM

package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/korovindenis/go-pc-metrics/internal/adapters/storage/memory/mocks"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewMemoryStorage(t *testing.T) {
	cfg := mocks.NewCfg(t)
	log := mocks.NewLog(t)
	log.On("Info", mock.Anything).Return("")
	storage, err := New(cfg, log)

	// Assert
	assert.NoError(t, err, "New should not return an error")
	assert.NotNil(t, storage, "Storage should not be nil")

}
func TestStorage_SaveAllData(t *testing.T) {
	cfg := mocks.NewCfg(t)
	log := mocks.NewLog(t)
	log.On("Info", mock.Anything).Return("")
	memory, _ := New(cfg, log)

	ctx := context.Background()
	floatVal := float64(1)
	intVal := int64(1)

	tests := []struct {
		name string
		args []entity.Metrics
		err  error
	}{
		{
			name: "positive gauge",
			args: []entity.Metrics{
				{
					ID:    "cpu",
					MType: "gauge",
					Value: &floatVal,
				},
			},
		},
		{
			name: "positive counter",
			args: []entity.Metrics{
				{
					ID:    "cpu",
					MType: "counter",
					Delta: &intVal,
				},
			},
		},
		{
			name: "negative",
			err:  errors.New("sendMetrics(): metricsVal not recognized"),
			args: []entity.Metrics{
				{
					ID:    "cpu",
					MType: "negative",
					Delta: &intVal,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := memory.SaveAllData(ctx, tt.args)

			// Assert
			if err != nil && tt.err == nil {
				t.Fatal(err)
			}
			if err != nil {
				if tt.args[0].MType == "gauge" {
					assert.Equal(t, memory.MetricsType.Gauge[tt.args[0].ID], *tt.args[0].Value)
				} else {
					assert.Equal(t, memory.MetricsType.Counter[tt.args[0].ID], *tt.args[0].Delta)
				}
			}
		})
	}
}

func TestStorage_GetAllData(t *testing.T) {
	cfg := mocks.NewCfg(t)
	log := mocks.NewLog(t)
	log.On("Info", mock.Anything).Return("")
	memory, _ := New(cfg, log)

	ctx := context.Background()

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "positive",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			data, err := memory.GetAllData(ctx)

			// Assert
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, data, memory.MetricsType)
		})
	}
}

func TestStorage_GetCounter(t *testing.T) {
	cfg := mocks.NewCfg(t)
	log := mocks.NewLog(t)
	log.On("Info", mock.Anything).Return("")
	memory, _ := New(cfg, log)

	ctx := context.Background()

	tests := []struct {
		name string
		arg  struct {
			name string
			val  int64
		}
		err error
	}{
		{
			name: "positive",
			arg: struct {
				name string
				val  int64
			}{
				name: "cpu",
				val:  int64(1),
			},
		},
		{
			name: "negative",
			err:  entity.ErrMetricNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Asset
			if tt.arg.name != "" {
				memory.MetricsType.Counter[tt.arg.name] = tt.arg.val
			}

			// Act
			data, err := memory.GetCounter(ctx, tt.arg.name)

			// Assert
			if err != tt.err {
				t.Fatal(err)
			}
			assert.Equal(t, data, memory.MetricsType.Counter[tt.arg.name])
		})
	}
}

func TestStorage_SaveCounter(t *testing.T) {
	cfg := mocks.NewCfg(t)
	log := mocks.NewLog(t)
	log.On("Info", mock.Anything).Return("")
	memory, _ := New(cfg, log)

	ctx := context.Background()

	tests := []struct {
		name string
		arg  struct {
			name string
			val  int64
		}
		err error
	}{
		{
			name: "positive",
			arg: struct {
				name string
				val  int64
			}{
				name: "cpu",
				val:  int64(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Asset
			memory.MetricsType.Counter[tt.arg.name] = tt.arg.val

			// Act
			err := memory.SaveCounter(ctx, tt.arg.name, tt.arg.val)

			// Assert
			if err != tt.err {
				t.Fatal(err)
			}
		})
	}
}

func TestStorage_GetGauge(t *testing.T) {
	cfg := mocks.NewCfg(t)
	log := mocks.NewLog(t)
	log.On("Info", mock.Anything).Return("")
	memory, _ := New(cfg, log)

	ctx := context.Background()

	tests := []struct {
		name string
		arg  struct {
			name string
			val  float64
		}
		err error
	}{
		{
			name: "positive",
			arg: struct {
				name string
				val  float64
			}{
				name: "cpu",
				val:  float64(1),
			},
		},
		{
			name: "negative",
			err:  entity.ErrMetricNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Asset
			if tt.arg.name != "" {
				memory.MetricsType.Gauge[tt.arg.name] = tt.arg.val
			}

			// Act
			data, err := memory.GetGauge(ctx, tt.arg.name)

			// Assert
			if err != tt.err {
				t.Fatal(err)
			}
			assert.Equal(t, data, memory.MetricsType.Gauge[tt.arg.name])
		})
	}
}

func TestStorage_SaveGauge(t *testing.T) {
	cfg := mocks.NewCfg(t)
	log := mocks.NewLog(t)
	log.On("Info", mock.Anything).Return("")
	memory, _ := New(cfg, log)

	ctx := context.Background()

	tests := []struct {
		name string
		arg  struct {
			name string
			val  float64
		}
		err error
	}{
		{
			name: "positive",
			arg: struct {
				name string
				val  float64
			}{
				name: "cpu",
				val:  float64(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Asset
			memory.MetricsType.Gauge[tt.arg.name] = tt.arg.val

			// Act
			err := memory.SaveGauge(ctx, tt.arg.name, tt.arg.val)

			// Assert
			if err != tt.err {
				t.Fatal(err)
			}
		})
	}
}

func TestStorage_Ping(t *testing.T) {
	log := mocks.NewLog(t)
	log.On("Info", mock.Anything).Return("")

	cfg := mocks.NewCfg(t)
	ctx := context.Background()
	memory, _ := New(cfg, log)

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "positive",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := memory.Ping(ctx)

			// Assert
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
