package agentapp

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
func Exec(agentUsecase agentUsecase, log logger, cfg config) error {
	restClient := resty.New()
	restClient.SetTimeout(time.Second * 10)

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
				return err
			}
			if err := agentUsecase.UpdateCounter(); err != nil {
				return err
			}
		case <-sendTicker.C:
			log.Info("send metrics")
			gaugeVal, err := agentUsecase.GetGauge()
			if err != nil {
				return err
			}
			err = sendMetrics(restClient, gaugeVal, log, httpServerAddress)
			if err != nil {
				return err
			}
			counterVal, err := agentUsecase.GetCounter()
			if err != nil {
				return err
			}
			err = sendMetrics(restClient, counterVal, log, httpServerAddress)
			if err != nil {
				return err
			}
		}
	}
}

// prepare data
func sendMetrics(restClient *resty.Client, metricsVal any, log logger, httpServerAddress string) error {
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
				return err
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
				return err
			}
		}
	default:
		return errors.New("sendMetrics(): metricsVal not recognized")
	}

	return nil
}

// send data
func httpReq(restClient *resty.Client, log logger, httpServerAddress string, metrics entity.Metrics) error {

	// payload, err := json.Marshal(metrics)
	// if err != nil {
	// 	return err
	// }

	//log.Info("Send: " + string(payload))

	_, err := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(metrics).
		Post(fmt.Sprintf("%s/update/", httpServerAddress))

	//HTTP POST request
	// req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/update/", httpServerAddress), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	// req.Close = true

	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Content-Length", strconv.Itoa(len(payload)))
	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()

	return nil
}
