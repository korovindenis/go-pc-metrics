package app

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

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
	if err := httpReq(restClient, log, httpServerAddress, metrics, 1000000); err != nil {
		return fmt.Errorf("sendMetrics entity.CounterType: %s", err)
	}
	return nil
}

// send data
func httpReq(restClient *http.Client, log logger, httpServerAddress string, metrics []entity.Metrics, batchSize int) error {
	// Проверяем, что размер батча больше нуля
	if batchSize <= 0 {
		return fmt.Errorf("batch size must be greater than zero")
	}

	// Создаем срез для хранения батчей
	batches := make([][]entity.Metrics, 0)

	// Разбиваем данные на батчи
	for i := 0; i < len(metrics); i += batchSize {
		end := i + batchSize
		if end > len(metrics) {
			end = len(metrics)
		}
		batch := metrics[i:end]
		batches = append(batches, batch)
	}

	// Отправляем каждый батч
	for _, batch := range batches {
		// Сначала сжимаем данные в формат Gzip
		var compressedBody bytes.Buffer
		gz := gzip.NewWriter(&compressedBody)

		jsonBody, err := json.Marshal(batch)
		if err != nil {
			return fmt.Errorf("err in Marshal: %s", err)
		}

		_, err = gz.Write(jsonBody)
		if err != nil {
			return fmt.Errorf("err in gz Write: %s", err)
		}

		gz.Close()

		// Создаем запрос с сжатыми данными и устанавливаем заголовки
		req, err := http.NewRequest("POST", httpServerAddress+"/updates/", &compressedBody)
		if err != nil {
			return fmt.Errorf("err in NewRequest: %s", err)
		}

		// Устанавливаем заголовки для сжатия и JSON
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")

		// Отправляем запрос
		resp, err := restClient.Do(req)
		if err != nil {
			log.Info(fmt.Sprintf("err in httpclient: %s", err))
		} else {
			defer resp.Body.Close()
			log.Info("Status Code:" + resp.Status)
		}
	}

	return nil
}

// send data
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
