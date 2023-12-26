// the business logic of the backend

package serverusecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	"github.com/korovindenis/go-pc-metrics/internal/domain/usecases/server/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServer_Ping(t *testing.T) {
	cfg := mocks.NewCfg(t)
	storage := mocks.NewStorage(t)
	server, _ := New(storage, cfg)

	cfg.On("GetServerAddress").Return("localhost:8080").Maybe()
	cfg.On("GetStoreInterval").Return(time.Duration(1 * time.Second)).Maybe()

	tests := []struct {
		name string
		err  error
		ctx  context.Context
	}{
		{
			name: "positive",
			ctx:  context.Background(),
		},
		{
			name: "negative",
			ctx:  context.Background(),
			err:  errors.New("err"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ping := storage.On("Ping", mock.Anything).Return(tt.err)

			// Act
			err := server.Ping(tt.ctx)

			// Assert
			assert.Equal(t, tt.err, err)

			// Unset
			ping.Unset()
		})
	}
}

func TestServer_SaveAllDataBatchUsecase(t *testing.T) {
	cfg := mocks.NewCfg(t)
	storage := mocks.NewStorage(t)
	server, _ := New(storage, cfg)

	cfg.On("GetServerAddress").Return("localhost:8080").Maybe()
	cfg.On("GetStoreInterval").Return(time.Duration(1 * time.Second)).Maybe()

	tests := []struct {
		name string
		arg  []entity.Metrics
		err  error
		ctx  context.Context
	}{
		{
			name: "positive",
			ctx:  context.Background(),
			arg:  []entity.Metrics{},
		},
		{
			name: "negative",
			ctx:  context.Background(),
			arg:  []entity.Metrics{},
			err:  errors.New("err"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			saveAllData := storage.On("SaveAllData", mock.Anything, mock.Anything).Return(tt.err)

			// Act
			err := server.SaveAllDataBatchUsecase(tt.ctx, tt.arg)

			// Assert
			assert.Equal(t, tt.err, err)

			// Unset
			saveAllData.Unset()
		})
	}
}

func TestServer_GetAllDataUsecase(t *testing.T) {
	cfg := mocks.NewCfg(t)
	storage := mocks.NewStorage(t)
	server, _ := New(storage, cfg)
	cfg.On("GetServerAddress").Return("localhost:8080").Maybe()
	cfg.On("GetStoreInterval").Return(time.Duration(1 * time.Second)).Maybe()

	tests := []struct {
		name string
		err  error
		ctx  context.Context
	}{
		{
			name: "positive",
			ctx:  context.Background(),
		},
		{
			name: "negative",
			err:  errors.New("err"),
			ctx:  context.Background(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			getAllData := storage.On("GetAllData", mock.Anything, mock.Anything).Return(entity.MetricsType{}, tt.err)

			// Act
			_, err := server.GetAllDataUsecase(tt.ctx)

			// Assert
			assert.Equal(t, tt.err, err)

			// Unset
			getAllData.Unset()
		})
	}
}

func TestServer_GetCounterUsecase(t *testing.T) {
	cfg := mocks.NewCfg(t)
	storage := mocks.NewStorage(t)
	server, _ := New(storage, cfg)
	cfg.On("GetServerAddress").Return("localhost:8080").Maybe()
	cfg.On("GetStoreInterval").Return(time.Duration(1 * time.Second)).Maybe()

	tests := []struct {
		name string
		err  error
		ctx  context.Context
	}{
		{
			name: "positive",
			ctx:  context.Background(),
		},
		{
			name: "negative",
			err:  errors.New("err"),
			ctx:  context.Background(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			getCounter := storage.On("GetCounter", mock.Anything, mock.Anything).Return(int64(0), tt.err)

			// Act
			_, err := server.GetCounterUsecase(tt.ctx, "cpu")

			// Assert
			assert.Equal(t, tt.err, err)

			// Unset
			getCounter.Unset()
		})
	}
}

func TestServer_GetGaugeUsecase(t *testing.T) {
	cfg := mocks.NewCfg(t)
	storage := mocks.NewStorage(t)
	server, _ := New(storage, cfg)
	cfg.On("GetServerAddress").Return("localhost:8080").Maybe()
	cfg.On("GetStoreInterval").Return(time.Duration(1 * time.Second)).Maybe()

	tests := []struct {
		name string
		err  error
		ctx  context.Context
	}{
		{
			name: "positive",
			ctx:  context.Background(),
		},
		{
			name: "negative",
			err:  errors.New("err"),
			ctx:  context.Background(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			getGauge := storage.On("GetGauge", mock.Anything, mock.Anything).Return(float64(0), tt.err)

			// Act
			_, err := server.GetGaugeUsecase(tt.ctx, "cpu")

			// Assert
			assert.Equal(t, tt.err, err)

			// Unset
			getGauge.Unset()
		})
	}
}

func TestServer_SaveGaugeUsecase(t *testing.T) {
	cfg := mocks.NewCfg(t)
	storage := mocks.NewStorage(t)
	server, _ := New(storage, cfg)
	cfg.On("GetServerAddress").Return("localhost:8080").Maybe()
	cfg.On("GetStoreInterval").Return(time.Duration(1 * time.Second)).Maybe()

	tests := []struct {
		name string
		err  error
		ctx  context.Context
	}{
		{
			name: "positive",
			ctx:  context.Background(),
		},
		{
			name: "negative",
			err:  errors.New("err"),
			ctx:  context.Background(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			getGauge := storage.On("SaveGauge", mock.Anything, mock.Anything, mock.Anything).Return(tt.err)

			// Act
			err := server.SaveGaugeUsecase(tt.ctx, "cpu", float64(0))

			// Assert
			assert.Equal(t, tt.err, err)

			// Unset
			getGauge.Unset()
		})
	}
}

func TestServer_SaveCounterUsecase(t *testing.T) {
	cfg := mocks.NewCfg(t)
	storage := mocks.NewStorage(t)
	server, _ := New(storage, cfg)
	cfg.On("GetServerAddress").Return("localhost:8080").Maybe()
	cfg.On("GetStoreInterval").Return(time.Duration(1 * time.Second)).Maybe()

	tests := []struct {
		name string
		err  error
		ctx  context.Context
	}{
		{
			name: "positive",
			ctx:  context.Background(),
		},
		{
			name: "negative",
			err:  errors.New("err"),
			ctx:  context.Background(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			getCounterUsecase := storage.On("SaveCounter", mock.Anything, mock.Anything, mock.Anything).Return(tt.err)
			getCounter := storage.On("GetCounter", mock.Anything, mock.Anything).Return(int64(0), tt.err)

			// Act
			err := server.SaveCounterUsecase(tt.ctx, "cpu", int64(0))

			// Assert
			assert.Equal(t, tt.err, err)

			// Unset
			getCounterUsecase.Unset()
			getCounter.Unset()
		})
	}
}

func TestServer_SaveAllDataUsecase(t *testing.T) {
	cfg := mocks.NewCfg(t)
	storage := mocks.NewStorage(t)
	server, _ := New(storage, cfg)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Local().Add(3*time.Second))

	cfg.On("GetServerAddress").Return("localhost:8080").Maybe()
	cfg.On("GetStoreInterval").Return(time.Duration(1 * time.Second)).Maybe()

	tests := []struct {
		name string
		arg  []entity.Metrics
		err  error
		ctx  context.Context
	}{
		{
			name: "positive",
			arg:  []entity.Metrics{},
			ctx:  ctx,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			saveAllData := storage.On("SaveAllData", mock.Anything, mock.Anything).Return(tt.err)

			// Act
			err := server.SaveAllDataUsecase(tt.ctx, tt.arg)

			// Assert
			assert.Equal(t, tt.err, err)

			// Unset
			saveAllData.Unset()

		})
	}
	cancel()
}
