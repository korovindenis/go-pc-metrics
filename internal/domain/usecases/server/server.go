package serverusecase

import (
	"context"
	"time"

	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
)

// storage functions
type storage interface {
	SaveGauge(ctx context.Context, gaugeName string, gaugeValue float64) error
	GetGauge(ctx context.Context, gaugeName string) (float64, error)

	SaveCounter(ctx context.Context, counterName string, counterValue int64) error
	GetCounter(ctx context.Context, counterName string) (int64, error)

	GetAllData(ctx context.Context) (entity.MetricsType, error)
	SaveAllData(ctx context.Context, metrics []entity.Metrics) error

	Ping(ctx context.Context) error
}

type cfg interface {
	GetServerAddress() string
	GetStoreInterval() time.Duration
}

type Server struct {
	storage       storage
	storeInterval time.Duration
}

func New(s any, config cfg) (*Server, error) {
	storageInstance, ok := s.(storage)
	if !ok {
		return nil, entity.ErrStorageInstance
	}

	return &Server{
		storage:       storageInstance,
		storeInterval: config.GetStoreInterval(),
	}, nil
}

func (s *Server) SaveGaugeUsecase(ctx context.Context, gaugeName string, gaugeValue float64) error {
	return s.storage.SaveGauge(ctx, gaugeName, gaugeValue)
}

func (s *Server) GetGaugeUsecase(ctx context.Context, gaugeName string) (float64, error) {
	return s.storage.GetGauge(ctx, gaugeName)
}

func (s *Server) SaveCounterUsecase(ctx context.Context, counterName string, counterValue int64) error {
	// current val + newVal
	currentCounterValue, err := s.GetCounterUsecase(ctx, counterName)
	if err != nil && err != entity.ErrMetricNotFound {
		return err
	}

	return s.storage.SaveCounter(ctx, counterName, counterValue+currentCounterValue)
}

func (s *Server) GetCounterUsecase(ctx context.Context, counterName string) (int64, error) {
	return s.storage.GetCounter(ctx, counterName)
}

func (s *Server) GetAllDataUsecase(ctx context.Context) (entity.MetricsType, error) {
	return s.storage.GetAllData(ctx)
}

func (s *Server) SaveAllDataUsecase(ctx context.Context, metrics []entity.Metrics) error {
	sendTicker := time.NewTicker(s.storeInterval)
	defer sendTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			// context is end, exit from for
			return nil
		case <-sendTicker.C:
			s.storage.SaveAllData(ctx, metrics)
		}
	}
}

func (s *Server) SaveAllDataBatchUsecase(ctx context.Context, metrics []entity.Metrics) error {
	var sumCounter int64
	for _, val := range metrics {
		if val.MType == "counter" {
			sumCounter += int64(*val.Delta)
		}
	}
	for key, val := range metrics {
		if val.MType == "counter" {
			metrics[key].Delta = &sumCounter
		}
	}

	return s.storage.SaveAllData(ctx, metrics)
}

func (s *Server) Ping(ctx context.Context) error {
	return s.storage.Ping(ctx)
}
