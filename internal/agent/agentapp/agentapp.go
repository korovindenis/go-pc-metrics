package agentapp

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
func Exec(agentUsecase agentUsecase, log logger, cfg config) error {
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
			err = sendMetrics(gaugeVal, log, httpServerAddress)
			if err != nil {
				return err
			}
			counterVal, err := agentUsecase.GetCounter()
			if err != nil {
				return err
			}
			err = sendMetrics(counterVal, log, httpServerAddress)
			if err != nil {
				return err
			}
		}
	}
}

// prepare data
func sendMetrics(metricsVal any, log logger, httpServerAddress string) error {
	switch v := metricsVal.(type) {
	case entity.GaugeType:
		_ = v // for go vet
		for name, value := range metricsVal.(entity.GaugeType) {
			if err := httpReq(log, httpServerAddress, "gauge", name, strconv.FormatFloat(value, 'g', -1, 64)); err != nil {
				return err
			}
		}
	case entity.CounterType:
		for name, value := range metricsVal.(entity.CounterType) {
			if err := httpReq(log, httpServerAddress, "counter", name, strconv.FormatInt(value, 10)); err != nil {
				return err
			}
		}
	default:
		return errors.New("sendMetrics(): metricsVal not recognized")
	}

	return nil
}

// send data
func httpReq(log logger, httpServerAddress, metricType, metricName, metricVal string) error {
	uri := fmt.Sprintf("%s/update/%s/%s/%s", httpServerAddress, metricType, metricName, metricVal)
	log.Info(uri)

	// HTTP POST request
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(""))
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
