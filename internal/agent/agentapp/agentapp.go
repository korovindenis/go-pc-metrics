package agentapp

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
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
	Println(v ...interface{})
}

// config functions
type config interface {
	GetHTTPAddressWithScheme() string
	GetPollInterval() time.Duration
	GetReportInterval() time.Duration
}

// agent main
func Exec(agentUsecase agentUsecase, log logger, cfg config) error {
	httpAddress := cfg.GetHTTPAddressWithScheme()
	reportInterval := cfg.GetReportInterval()
	pollInterval := cfg.GetPollInterval()

	for mainTime := reportInterval; ; mainTime -= pollInterval {
		// every 2 sec.
		log.Println("update metrics")
		if err := agentUsecase.UpdateGauge(); err != nil {
			return err
		}
		if err := agentUsecase.UpdateCounter(); err != nil {
			return err
		}

		// every 10 sec.
		if mainTime <= 0 {
			log.Println("send metrics")
			gaugeVal, err := agentUsecase.GetGauge()
			if err != nil {
				return err
			}
			err = sendMetrics(gaugeVal, log, httpAddress)
			if err != nil {
				return err
			}

			counterVal, err := agentUsecase.GetCounter()
			if err != nil {
				return err
			}
			err = sendMetrics(counterVal, log, httpAddress)
			if err != nil {
				return err
			}
			mainTime = reportInterval
		}

		time.Sleep(pollInterval)
	}
}

// prepare data
func sendMetrics(metricsVal any, log logger, httpAddress string) error {
	switch v := metricsVal.(type) {
	case entity.GaugeType:
		_ = v
		for name, value := range metricsVal.(entity.GaugeType) {
			if err := httpReq(log, httpAddress, "gauge", name, strconv.FormatFloat(value, 'g', -1, 64)); err != nil {
				return err
			}
		}
	case entity.CounterType:
		for name, value := range metricsVal.(entity.CounterType) {
			if err := httpReq(log, httpAddress, "counter", name, strconv.FormatInt(value, 10)); err != nil {
				return err
			}
		}
	default:
		return errors.New("sendMetrics(): metricsVal not recognized")
	}

	return nil
}

// send data
func httpReq(log logger, httpAddress, metricType, metricName, metricVal string) error {
	uri := fmt.Sprintf("%s/update/%s/%s/%s", httpAddress, metricType, metricName, metricVal)
	log.Println(uri)

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
