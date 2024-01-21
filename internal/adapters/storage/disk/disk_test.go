// Storage with file system

package disk

import (
	"context"
	"testing"

	"github.com/korovindenis/go-pc-metrics/internal/adapters/storage/disk/mocks"
	"github.com/stretchr/testify/mock"
)

func TestStorage_SaveGauge(t *testing.T) {
	cfg := mocks.NewCfg(t)
	cfg.On("GetFileStoragePath").Return("").Maybe()
	cfg.On("GetRestore").Return(false).Maybe()

	log := mocks.NewLog(t)
	log.On("Info", mock.Anything).Return("")

	disk, _ := New(cfg, log)

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
			disk.metrics.Gauge[tt.arg.name] = tt.arg.val

			// Act
			err := disk.SaveGauge(ctx, tt.arg.name, tt.arg.val)

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
	cfg.On("GetFileStoragePath").Return("").Maybe()
	cfg.On("GetRestore").Return(false).Maybe()

	ctx := context.Background()
	disk, _ := New(cfg, log)

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
			err := disk.Ping(ctx)

			// Assert
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
