package agentusecase

import (
	"math/rand"
	"runtime"

	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
)

type Agent struct {
	runtime runtime.MemStats
	metrics entity.MetricsType
}

func New() (*Agent, error) {
	agentUsecase := &Agent{
		metrics: entity.MetricsType{
			Gauge:   make(map[string]float64, 30),
			Counter: make(map[string]int64, 1),
		},
	}
	runtime.ReadMemStats(&agentUsecase.runtime)

	return agentUsecase, nil
}

func (a *Agent) UpdateCounter() error {
	a.metrics.Counter["PollCount"] += 1

	return nil
}

func (a *Agent) UpdateGauge() error {

	a.metrics.Gauge = entity.GaugeType{
		"Alloc":         rand.Float64(), //float64(a.runtime.Alloc),
		"BuckHashSys":   float64(a.runtime.BuckHashSys),
		"Frees":         rand.Float64(), //float64(a.runtime.Frees),
		"GCCPUFraction": a.runtime.GCCPUFraction,
		"GCSys":         float64(a.runtime.GCSys),
		"HeapAlloc":     rand.Float64(), //float64(a.runtime.HeapAlloc),
		"HeapIdle":      rand.Float64(), //float64(a.runtime.HeapIdle),
		"HeapInuse":     rand.Float64(), //float64(a.runtime.HeapInuse),
		"HeapObjects":   rand.Float64(), //float64(a.runtime.HeapObjects),
		"HeapReleased":  float64(a.runtime.HeapReleased),
		"HeapSys":       float64(a.runtime.HeapSys),
		"LastGC":        float64(a.runtime.LastGC),
		"Lookups":       float64(a.runtime.Lookups),
		"MCacheSys":     rand.Float64(), //float64(a.runtime.MCacheSys),
		"MSpanInuse":    rand.Float64(), //float64(a.runtime.MSpanInuse),
		"MCacheInuse":   rand.Float64(), //float64(a.runtime.MSpanInuse),
		"MSpanSys":      float64(a.runtime.MSpanSys),
		"Mallocs":       rand.Float64(), // float64(a.runtime.Mallocs),
		"NextGC":        float64(a.runtime.NextGC),
		"NumForcedGC":   float64(a.runtime.NumForcedGC),
		"NumGC":         float64(a.runtime.NumGC),
		"OtherSys":      float64(a.runtime.OtherSys),
		"PauseTotalNs":  float64(a.runtime.PauseTotalNs),
		"StackInuse":    float64(a.runtime.StackInuse),
		"StackSys":      float64(a.runtime.StackSys),
		"Sys":           float64(a.runtime.Sys),
		"TotalAlloc":    rand.Float64(), // float64(a.runtime.TotalAlloc),
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
