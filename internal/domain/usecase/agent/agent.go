package agentusecase

import (
	"math/rand"
	"runtime"

	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
)

// agent functions
type AgentUsecase interface {
	GetGauge() (entity.GaugeType, error)
	GetCounter() (entity.CounterType, error)

	UpdateGauge() error
	UpdateCounter() error
}

type Agent struct {
	rtm     runtime.MemStats
	metrics entity.MetricsType
}

func New() (*Agent, error) {
	agntUscs := &Agent{
		metrics: entity.MetricsType{
			Gauge:   make(map[string]float64, 30),
			Counter: make(map[string]int64, 1),
		},
	}
	runtime.ReadMemStats(&agntUscs.rtm)

	return agntUscs, nil
}

func (a *Agent) UpdateCounter() error {
	a.metrics.Counter["PollCount"] += 1

	return nil
}

func (a *Agent) UpdateGauge() error {

	a.metrics.Gauge = entity.GaugeType{
		"Alloc":         float64(a.rtm.Alloc),
		"BuckHashSys":   float64(a.rtm.BuckHashSys),
		"Frees":         float64(a.rtm.Frees),
		"GCCPUFraction": a.rtm.GCCPUFraction,
		"GCSys":         float64(a.rtm.GCSys),
		"HeapAlloc":     float64(a.rtm.HeapAlloc),
		"HeapIdle":      float64(a.rtm.HeapIdle),
		"HeapInuse":     float64(a.rtm.HeapInuse),
		"HeapObjects":   float64(a.rtm.HeapObjects),
		"HeapReleased":  float64(a.rtm.HeapReleased),
		"HeapSys":       float64(a.rtm.HeapSys),
		"LastGC":        float64(a.rtm.LastGC),
		"Lookups":       float64(a.rtm.Lookups),
		"MCacheSys":     float64(a.rtm.MCacheSys),
		"MSpanInuse":    float64(a.rtm.MSpanInuse),
		"MSpanSys":      float64(a.rtm.MSpanSys),
		"Mallocs":       float64(a.rtm.Mallocs),
		"NextGC":        float64(a.rtm.NextGC),
		"NumForcedGC":   float64(a.rtm.NumForcedGC),
		"NumGC":         float64(a.rtm.NumGC),
		"OtherSys":      float64(a.rtm.OtherSys),
		"PauseTotalNs":  float64(a.rtm.PauseTotalNs),
		"StackInuse":    float64(a.rtm.StackInuse),
		"StackSys":      float64(a.rtm.StackSys),
		"Sys":           float64(a.rtm.Sys),
		"TotalAlloc":    float64(a.rtm.TotalAlloc),
		"RandomValue":   rand.Float64(),
	}

	return nil
}

func (a *Agent) GetGauge() (entity.GaugeType, error) {
	return a.metrics.Gauge, nil
}
func (a *Agent) GetCounter() (entity.CounterType, error) {
	return a.metrics.Counter, nil
}
