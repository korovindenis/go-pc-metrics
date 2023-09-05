package serverusecase

import (
	"time"

	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
)

// storage functions
type storage interface {
	SaveGauge(gaugeName string, gaugeValue float64) error
	GetGauge(gaugeName string) (float64, error)

	SaveCounter(counterName string, counterValue int64) error
	GetCounter(counterName string) (int64, error)

	GetAllData() (entity.MetricsType, error)
	SaveAllData() error
}

type cfg interface {
	GetServerAddress() string
	GetStoreInterval() time.Duration
}

type Server struct {
	storage storage
}

func New(s storage) (*Server, error) {
	return &Server{
		storage: s,
	}, nil
}

func (s *Server) SaveGaugeUsecase(gaugeName string, gaugeValue float64) error {
	return s.storage.SaveGauge(gaugeName, gaugeValue)
}

func (s *Server) GetGaugeUsecase(gaugeName string) (float64, error) {
	return s.storage.GetGauge(gaugeName)
}

func (s *Server) SaveCounterUsecase(counterName string, counterValue int64) error {
	// current val + newVal
	currentCounterValue, err := s.GetCounterUsecase(counterName)
	if err != nil && err != entity.ErrMetricNotFound {
		return err
	}

	return s.storage.SaveCounter(counterName, counterValue+currentCounterValue)
}

func (s *Server) GetCounterUsecase(counterName string) (int64, error) {
	return s.storage.GetCounter(counterName)
}

func (s *Server) GetAllDataUsecase() (entity.MetricsType, error) {
	return s.storage.GetAllData()
}

func (s *Server) SaveAllDataUsecase(cfg cfg) {
	sendTicker := time.NewTicker(cfg.GetStoreInterval())
	defer sendTicker.Stop()

	for range sendTicker.C {
		s.storage.SaveAllData()
	}
}
