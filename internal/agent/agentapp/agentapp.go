package agentapp

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	"go.uber.org/zap/zapcore"
)

// agent functions
type agentUsecase interface {
	GetGauge() (entity.GaugeType, error)
	GetCounter() (entity.CounterType, error)

	UpdateGauge() error
	UpdateCounter() error
}

// logger functions
type logger interface {
	Info(msg string, fields ...zapcore.Field)
}

// config functions
type config interface {
	GetServerAddressWithScheme() string
	GetPollInterval() time.Duration
	GetReportInterval() time.Duration
}

// agent main
func Exec(agentUsecase agentUsecase, log logger, cfg config) error {
	restClient := &http.Client{}

	httpServerAddress := cfg.GetServerAddressWithScheme()

	updateTicker := time.NewTicker(cfg.GetPollInterval())
	defer updateTicker.Stop()

	sendTicker := time.NewTicker(cfg.GetReportInterval())
	defer sendTicker.Stop()

	for {
		select {
		case <-updateTicker.C:
			log.Info("update metrics")
			if err := agentUsecase.UpdateGauge(); err != nil {
				return fmt.Errorf("agentapp Exec UpdateGauge: %s", err)
			}
			if err := agentUsecase.UpdateCounter(); err != nil {
				return fmt.Errorf("agentapp Exec UpdateCounter: %s", err)
			}
		case <-sendTicker.C:
			log.Info("send metrics")
			gaugeVal, err := agentUsecase.GetGauge()
			if err != nil {
				return fmt.Errorf("agentapp Exec GetGauge: %s", err)
			}
			err = sendMetrics(restClient, gaugeVal, log, httpServerAddress)
			if err != nil {
				return fmt.Errorf("agentapp Exec sendMetrics: %s", err)
			}
			counterVal, err := agentUsecase.GetCounter()
			if err != nil {
				return fmt.Errorf("agentapp Exec GetCounter: %s", err)
			}
			err = sendMetrics(restClient, counterVal, log, httpServerAddress)
			if err != nil {
				return fmt.Errorf("agentapp Exec sendMetrics: %s", err)
			}
		}
	}
}

// prepare data
func sendMetrics(restClient *http.Client, metricsVal any, log logger, httpServerAddress string) error {
	switch v := metricsVal.(type) {
	case entity.GaugeType:
		_ = v // for go vet
		for name, value := range metricsVal.(entity.GaugeType) {
			metrics := entity.Metrics{
				ID:    name,
				MType: "gauge",
				Value: &value,
			}
			if err := httpReq(restClient, log, httpServerAddress, metrics); err != nil {
				return fmt.Errorf("sendMetrics entity.GaugeType: %s", err)
			}
		}
	case entity.CounterType:
		for name, value := range metricsVal.(entity.CounterType) {
			metrics := entity.Metrics{
				ID:    name,
				MType: "counter",
				Delta: &value,
			}
			if err := httpReq(restClient, log, httpServerAddress, metrics); err != nil {
				return fmt.Errorf("sendMetrics entity.CounterType: %s", err)
			}
		}
	default:
		return errors.New("sendMetrics(): metricsVal not recognized")
	}

	return nil
}

// send data
func httpReq(restClient *http.Client, log logger, httpServerAddress string, metrics entity.Metrics) error {

	// resty
	resty.New().SetDebug(true).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(&metrics).
		Post(httpServerAddress + "/update/")

	//if err != nil {
	//	return fmt.Errorf("err in resty: %s", err)
	//}

	return nil
}
