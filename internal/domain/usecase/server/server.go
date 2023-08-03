package serverUsecase

import "github.com/korovindenis/go-pc-metrics/internal/domain/entity"

type serverUsecase struct {
	storage entity.IStorage
}

type IServerUsecase interface {
	SaveGauge(gaugeName string, gaugeValue float64) error
	SaveCounter(counterName string, counterValue int64) error
}

func New(storage entity.IStorage) (IServerUsecase, error) {
	return &serverUsecase{
		storage: storage,
	}, nil
}

func (s *serverUsecase) SaveGauge(gaugeName string, gaugeValue float64) error {
	return s.storage.SaveGauge(gaugeName, gaugeValue)
}

func (s *serverUsecase) SaveCounter(counterName string, counterValue int64) error {
	// current val + newVal
	currentCounterValue, err := s.storage.GetCounter(counterName)
	if err != nil {
		return err
	}

	return s.storage.SaveCounter(counterName, counterValue+currentCounterValue)
}
