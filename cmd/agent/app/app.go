package app

import (
	"errors"
	"fmt"
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
	// Создаем новый запрос с использованием Resty
	resp, err := restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(metrics).
		EnableTrace().
		Post(httpServerAddress + "/updates/")

	if err != nil {
		log.Info(fmt.Sprintf("err in httpclient: %s", err))
		return err
	}

	// Проверяем код статуса ответа
	if resp.IsError() {
		log.Info("Status Code:" + resp.Status())
		log.Info("HTTP Error: %s" + resp.Status())
		log.Info("Response Body: %s" + resp.String())
	}

	// Если нужно, можно получить ответ
	// responseBody := resp.Body()

	return nil
}

// // send data
// func httpReq(restClient *http.Client, log logger, httpServerAddress string, metrics []entity.Metrics) error {
// 	var compressedBody bytes.Buffer
// 	gz := gzip.NewWriter(&compressedBody)

// 	jsonBody, err := json.Marshal(metrics)
// 	if err != nil {
// 		return fmt.Errorf("err in Marshal: %s", err)
// 	}

// 	log.Info("send: " + string(jsonBody))

// 	_, err = gz.Write(jsonBody)
// 	if err != nil {
// 		return fmt.Errorf("err in gz Write: %s", err)
// 	}
// 	gz.Close()

// 	req, err := http.NewRequest("POST", httpServerAddress+"/updates/", &compressedBody)
// 	if err != nil {
// 		return fmt.Errorf("err in NewRequest: %s", err)
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Content-Encoding", "gzip")
// 	req.Header.Set("Accept-Encoding", "gzip")

// 	resp, err := restClient.Do(req)
// 	if err != nil {
// 		log.Info(fmt.Sprintf("err in httpclient: %s", err))
// 	} else {
// 		defer resp.Body.Close()
// 		log.Info("Status Code:" + resp.Status)
// 	}

// 	return nil
// }
