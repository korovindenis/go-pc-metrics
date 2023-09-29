package app

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	GetKey() string
	GetRateLimit() int
}

type resultWorkerMetric struct {
	data bool
	err  error
}

// agent main
func Run(ctx context.Context, agentUsecase agentUsecase, log logger, cfg config) error {
	rateLimitCh := make(chan bool, cfg.GetRateLimit())
	resultCh := make(chan resultWorkerMetric)

	go updateWorker(ctx, agentUsecase, log, cfg, resultCh)
	go sendWorker(ctx, agentUsecase, log, cfg, resultCh, rateLimitCh)

	for res := range resultCh {
		if res.err != nil {
			return fmt.Errorf("agentapp Exec updateWorker: %s", res.err)
		}
	}

	return nil
}

func updateWorker(ctx context.Context, agentUsecase agentUsecase, log logger, cfg config, resultCh chan resultWorkerMetric) {
	updateTicker := time.NewTicker(cfg.GetPollInterval())
	defer updateTicker.Stop()
	defer close(resultCh)

	for {
		select {
		case <-ctx.Done():
			return
		case <-updateTicker.C:
			log.Info("update metrics")
			if err := agentUsecase.UpdateGauge(); err != nil {
				resultCh <- resultWorkerMetric{
					err: err,
				}
			}
			if err := agentUsecase.UpdateCounter(); err != nil {
				resultCh <- resultWorkerMetric{
					err: err,
				}
			}
			resultCh <- resultWorkerMetric{
				data: true,
			}
		}
	}
}

func sendWorker(ctx context.Context, agentUsecase agentUsecase, log logger, cfg config, resultCh chan resultWorkerMetric, rateLimitCh chan bool) {
	restClient := resty.New()
	httpServerAddress := cfg.GetServerAddressWithScheme()
	sendTicker := time.NewTicker(cfg.GetReportInterval())
	secretKey := cfg.GetKey()
	defer sendTicker.Stop()
	defer close(resultCh)

	for {
		select {
		case <-ctx.Done():
			return
		case <-sendTicker.C:
			// for limit request
			rateLimitCh <- true

			log.Info("send metrics")
			gaugeVal, err := agentUsecase.GetGauge()
			if err != nil {
				resultCh <- resultWorkerMetric{
					err: err,
				}
			}
			err = sendMetrics(restClient, gaugeVal, log, httpServerAddress, secretKey)
			if err != nil {
				resultCh <- resultWorkerMetric{
					err: err,
				}
			}
			counterVal, err := agentUsecase.GetCounter()
			if err != nil {
				resultCh <- resultWorkerMetric{
					err: err,
				}
			}
			err = sendMetrics(restClient, counterVal, log, httpServerAddress, secretKey)
			if err != nil {
				resultCh <- resultWorkerMetric{
					err: err,
				}
			}
			resultCh <- resultWorkerMetric{
				data: true,
			}

			<-rateLimitCh
		}
	}
}

// prepare data
func sendMetrics(restClient *resty.Client, metricsVal any, log logger, httpServerAddress, secretKey string) error {
	var metrics []entity.Metrics

	switch v := metricsVal.(type) {
	case entity.GaugeType:
		_ = v // for go vet
		for name, value := range metricsVal.(entity.GaugeType) {
			floatValue := new(float64)
			*floatValue = value
			metrics = append(metrics, entity.Metrics{
				ID:    name,
				MType: "gauge",
				Value: floatValue,
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

	if err := httpReq(restClient, log, httpServerAddress, secretKey, metrics); err != nil {
		return fmt.Errorf("sendMetrics entity.CounterType: %s", err)
	}
	return nil
}

// send data
func httpReq(restyClient *resty.Client, log logger, httpServerAddress, secretKey string, metrics []entity.Metrics) error {

	jsonBody, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("error in Marshal: %s", err)
	}

	log.Info("Send Metrics: " + string(jsonBody))

	var compressedBody bytes.Buffer
	gz := gzip.NewWriter(&compressedBody)
	_, err = gz.Write(jsonBody)
	if err != nil {
		return fmt.Errorf("error in gz Write: %s", err)
	}
	gz.Close()

	req := restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Length", strconv.Itoa(compressedBody.Len())).
		SetBody(compressedBody.Bytes()).
		EnableTrace()

	var hashSHA256 string
	if secretKey != "" {
		hashSHA256, _ = computeHMAC([]byte(jsonBody), secretKey)
		req.SetHeader("HashSHA256", hashSHA256)
		log.Info("HashSHA256: " + hashSHA256)
	}

	resp, err := req.Execute("POST", httpServerAddress+"/updates/")
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

func computeHMAC(input []byte, key string) (string, error) {
	keyBytes := []byte(key)

	h := hmac.New(sha256.New, keyBytes)

	_, err := h.Write(input)
	if err != nil {
		return "", err
	}

	hashBytes := h.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}
