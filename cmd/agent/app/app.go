package app

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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
func Run(agentUsecase agentUsecase, log logger, cfg config) error {
	restClient := resty.New()
	restClient.SetDebug(true)

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
func sendMetrics(restClient *resty.Client, metricsVal any, log logger, httpServerAddress string) error {
	var metrics []entity.Metrics

	switch v := metricsVal.(type) {
	case entity.GaugeType:
		_ = v // for go vet
		for name, value := range metricsVal.(entity.GaugeType) {
			metrics = append(metrics, entity.Metrics{
				ID:    name,
				MType: "gauge",
				Value: &value,
			})
		}
	case entity.CounterType:
		for name, value := range metricsVal.(entity.CounterType) {
			metrics = append(metrics, entity.Metrics{
				ID:    name,
				MType: "counter",
				Delta: &value,
			})
		}
	default:
		return errors.New("sendMetrics(): metricsVal not recognized")
	}
	if err := httpReq(restClient, log, httpServerAddress, metrics); err != nil {
		return fmt.Errorf("sendMetrics entity.CounterType: %s", err)
	}
	return nil
}

// send data
func httpReq(restyClient *resty.Client, log logger, httpServerAddress string, metrics []entity.Metrics) error {
	float := 42.5
	metrics = []entity.Metrics{
		entity.Metrics{
			ID:    "test_id",
			MType: "gauge",
			Value: &float,
		},
	}

	jsonBody, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("error in Marshal: %s", err)
	}

	fmt.Println("Send Metrics:", string(jsonBody))

	var compressedBody bytes.Buffer
	gz := gzip.NewWriter(&compressedBody)
	_, err = gz.Write(jsonBody)
	if err != nil {
		return fmt.Errorf("error in gz Write: %s", err)
	}
	gz.Close()

	resp, err := restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Length", strconv.Itoa(compressedBody.Len())).
		SetBody(compressedBody.Bytes()).
		EnableTrace().
		Post(httpServerAddress + "/updates/")

	if err != nil {
		log.Info(fmt.Sprintf("error in httpclient: %s", err))
	}

	if resp.IsError() {
		log.Info("Status Code:" + resp.Status())
		log.Info("HTTP Error: " + resp.Status())
		log.Info("Response Body: " + resp.String())
	}
	return nil
}
