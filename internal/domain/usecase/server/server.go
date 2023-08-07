package serverusecase

import "github.com/korovindenis/go-pc-metrics/internal/domain/entity"

// storage functions
type IStorage interface {
	SaveGauge(gaugeName string, gaugeValue float64) error
	GetGauge(gaugeName string) (float64, error)

	SaveCounter(counterName string, counterValue int64) error
	GetCounter(counterName string) (int64, error)

	GetAllData() (entity.MetricsType, error)
}

// server functions
type IServerUsecase interface {
	IStorage
}

type serverUsecase struct {
	storage IStorage
}

func New(storage IStorage) (IServerUsecase, error) {
	return &serverUsecase{
		storage: storage,
	}, nil
}

func (s *serverUsecase) SaveGauge(gaugeName string, gaugeValue float64) error {
	return s.storage.SaveGauge(gaugeName, gaugeValue)
}

func (s *serverUsecase) GetGauge(gaugeName string) (float64, error) {
	return s.storage.GetGauge(gaugeName)
}

func (s *serverUsecase) SaveCounter(counterName string, counterValue int64) error {
	// current val + newVal
	currentCounterValue, err := s.GetCounter(counterName)
	if err != nil && err != entity.MetricNotFoundErr {
		return err
	}

	return s.storage.SaveCounter(counterName, counterValue+currentCounterValue)
}

func (s *serverUsecase) GetCounter(counterName string) (int64, error) {
	return s.storage.GetCounter(counterName)
}

func (s *serverUsecase) GetAllData() (entity.MetricsType, error) {
	return s.storage.GetAllData()
}
