package agentusecase

import (
	"math/rand"
	"runtime"

	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
)

// agent functions
type IAgentUsecase interface {
	GetGauge() (entity.GaugeType, error)
	GetCounter() (entity.CounterType, error)

	UpdateGauge() error
	UpdateCounter() error
}

type agentUsecase struct {
	Rtm     runtime.MemStats
	Metrics entity.MetricsType
}

func New() (IAgentUsecase, error) {
	agntUscs := &agentUsecase{
		Metrics: entity.MetricsType{
			Gauge:   make(map[string]float64, 30),
			Counter: make(map[string]int64, 1),
		},
	}
	runtime.ReadMemStats(&agntUscs.Rtm)

	return agntUscs, nil
}

func (a *agentUsecase) UpdateCounter() error {
	a.Metrics.Counter["PollCount"] += 1

	return nil
}

func (a *agentUsecase) UpdateGauge() error {

	a.Metrics.Gauge = entity.GaugeType{
		"Alloc":         float64(a.Rtm.Alloc),
		"BuckHashSys":   float64(a.Rtm.BuckHashSys),
		"Frees":         float64(a.Rtm.Frees),
		"GCCPUFraction": a.Rtm.GCCPUFraction,
		"GCSys":         float64(a.Rtm.GCSys),
		"HeapAlloc":     float64(a.Rtm.HeapAlloc),
		"HeapIdle":      float64(a.Rtm.HeapIdle),
		"HeapInuse":     float64(a.Rtm.HeapInuse),
		"HeapObjects":   float64(a.Rtm.HeapObjects),
		"HeapReleased":  float64(a.Rtm.HeapReleased),
		"HeapSys":       float64(a.Rtm.HeapSys),
		"LastGC":        float64(a.Rtm.LastGC),
		"Lookups":       float64(a.Rtm.Lookups),
		"MCacheSys":     float64(a.Rtm.MCacheSys),
		"MSpanInuse":    float64(a.Rtm.MSpanInuse),
		"MSpanSys":      float64(a.Rtm.MSpanSys),
		"Mallocs":       float64(a.Rtm.Mallocs),
		"NextGC":        float64(a.Rtm.NextGC),
		"NumForcedGC":   float64(a.Rtm.NumForcedGC),
		"NumGC":         float64(a.Rtm.NumGC),
		"OtherSys":      float64(a.Rtm.OtherSys),
		"PauseTotalNs":  float64(a.Rtm.PauseTotalNs),
		"StackInuse":    float64(a.Rtm.StackInuse),
		"StackSys":      float64(a.Rtm.StackSys),
		"Sys":           float64(a.Rtm.Sys),
		"TotalAlloc":    float64(a.Rtm.TotalAlloc),
		"RandomValue":   rand.Float64(),
	}

	return nil
}

func (a *agentUsecase) GetGauge() (entity.GaugeType, error) {
	return a.Metrics.Gauge, nil
}
func (a *agentUsecase) GetCounter() (entity.CounterType, error) {
	return a.Metrics.Counter, nil
}
