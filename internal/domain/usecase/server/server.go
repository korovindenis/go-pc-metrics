package serverusecase

import "github.com/korovindenis/go-pc-metrics/internal/domain/entity"

// storage functions
type Storage interface {
	SaveGauge(gaugeName string, gaugeValue float64) error
	GetGauge(gaugeName string) (float64, error)

	SaveCounter(counterName string, counterValue int64) error
	GetCounter(counterName string) (int64, error)

	GetAllData() (entity.MetricsType, error)
}

// server functions
type ServerUsecase interface {
	Storage
}

type Server struct {
	storage Storage
}

func New(storage Storage) (*Server, error) {
	return &Server{
		storage: storage,
	}, nil
}

func (s *Server) SaveGauge(gaugeName string, gaugeValue float64) error {
	return s.storage.SaveGauge(gaugeName, gaugeValue)
}

func (s *Server) GetGauge(gaugeName string) (float64, error) {
	return s.storage.GetGauge(gaugeName)
}

func (s *Server) SaveCounter(counterName string, counterValue int64) error {
	// current val + newVal
	currentCounterValue, err := s.GetCounter(counterName)
	if err != nil && err != entity.ErrMetricNotFound {
		return err
	}

	return s.storage.SaveCounter(counterName, counterValue+currentCounterValue)
}

func (s *Server) GetCounter(counterName string) (int64, error) {
	return s.storage.GetCounter(counterName)
}

func (s *Server) GetAllData() (entity.MetricsType, error) {
	return s.storage.GetAllData()
}
