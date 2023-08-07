package agentapp

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/korovindenis/go-pc-metrics/internal/adapter/logger"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	agentusecase "github.com/korovindenis/go-pc-metrics/internal/domain/usecase/agent"
)

// agent main
func Exec(agntUscs agentusecase.IAgentUsecase, logger logger.ILogger, httpAddress string, pollInterval, reportInterval time.Duration) error {

	for mainTime := reportInterval; ; mainTime -= pollInterval {
		// every 2 sec.
		logger.Println("update metrics")
		if err := agntUscs.UpdateGauge(); err != nil {
			return err
		}
		if err := agntUscs.UpdateCounter(); err != nil {
			return err
		}

		// every 10 sec.
		if mainTime <= 0 {
			logger.Println("send metrics")
			// TODO err
			gaugeVal, _ := agntUscs.GetGauge()
			if err := sendMetrics(gaugeVal, logger, httpAddress); err != nil {
				return err
			}
			// TODO err
			counterVal, _ := agntUscs.GetCounter()
			if err := sendMetrics(counterVal, logger, httpAddress); err != nil {
				return err
			}
			mainTime = reportInterval
		}

		time.Sleep(pollInterval)
	}
}

// prepare data
func sendMetrics(metricsVal any, logger logger.ILogger, httpAddress string) error {
	switch v := metricsVal.(type) {
	case entity.GaugeType:
		_ = v
		for name, value := range metricsVal.(entity.GaugeType) {
			if err := httpReq(logger, httpAddress, "gauge", name, strconv.FormatFloat(value, 'f', 2, 64)); err != nil {
				return err
			}
		}
	case entity.CounterType:
		for name, value := range metricsVal.(entity.CounterType) {
			if err := httpReq(logger, httpAddress, "counter", name, strconv.FormatInt(value, 10)); err != nil {
				return err
			}
		}
	default:
		return errors.New("sendMetrics(): metricsVal not recognized")
	}

	return nil
}

// send data
func httpReq(logger logger.ILogger, httpAddress, metricType, metricName, metricVal string) error {
	uri := fmt.Sprintf("%s/update/%s/%s/%s", httpAddress, metricType, metricName, metricVal)
	logger.Println(uri)

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
