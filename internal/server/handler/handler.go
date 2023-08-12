package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
)

// function usecase
type usecase interface {
	SaveGaugeUsecase(gaugeName string, gaugeValue float64) error
	GetGaugeUsecase(gaugeName string) (float64, error)

	SaveCounterUsecase(counterName string, counterValue int64) error
	GetCounterUsecase(counterName string) (int64, error)

	GetAllDataUsecase() (entity.MetricsType, error)
}

type Handler struct {
	srvUsecase usecase
}

func New(u usecase) (*Handler, error) {
	return &Handler{
		srvUsecase: u,
	}, nil
}

func (s *Handler) ReceptionMetrics(c *gin.Context) {
	// get metric from url
	namedURL := entity.MetricsURI{
		MetricType: c.Param("metricType"),
		MetricName: c.Param("metricName"),
		MetricVal:  c.Param("metricVal"),
	}

	// validate metrics
	if namedURL.MetricType == "" || namedURL.MetricName == "" || namedURL.MetricVal == "" {
		c.AbortWithError(http.StatusNotFound, entity.ErrInvalidURLFormat)
		return
	}

	// run usecases
	switch namedURL.MetricType {
	case "gauge":
		// to float64
		metricVal, err := strconv.ParseFloat(namedURL.MetricVal, 64)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, entity.ErrInputVarIsWrongType)
			return
		}

		// save metric
		if err = s.srvUsecase.SaveGaugeUsecase(namedURL.MetricName, metricVal); err != nil {
			c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
			return
		}
	case "counter":
		// to int64
		metricVal, err := strconv.ParseInt(namedURL.MetricVal, 10, 64)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, entity.ErrInputVarIsWrongType)
			return
		}

		// save metric
		if err = s.srvUsecase.SaveCounterUsecase(namedURL.MetricName, metricVal); err != nil {
			c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
			return
		}
	default:
		c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Handler) OutputMetric(c *gin.Context) {
	// get metric from url
	namedURL := entity.MetricsURI{
		MetricType: c.Param("metricType"),
		MetricName: c.Param("metricName"),
	}

	// validate metrics
	if namedURL.MetricType == "" || namedURL.MetricName == "" {
		c.AbortWithError(http.StatusNotFound, entity.ErrInvalidURLFormat)
		return
	}

	switch namedURL.MetricType {
	case "gauge":
		// get metric
		gaugeVal, err := s.srvUsecase.GetGaugeUsecase(namedURL.MetricName)
		if err != nil {
			if errors.Is(err, entity.ErrMetricNotFound) {
				c.AbortWithError(http.StatusNotFound, entity.ErrInputMetricNotFound)
				return
			}
			c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
			return
		}

		// show metric
		c.String(http.StatusOK, strconv.FormatFloat(gaugeVal, 'g', -1, 64))
	case "counter":
		// get metric
		counterVal, err := s.srvUsecase.GetCounterUsecase(namedURL.MetricName)
		if err != nil {
			if errors.Is(err, entity.ErrMetricNotFound) {
				c.AbortWithError(http.StatusNotFound, entity.ErrInputMetricNotFound)
				return
			}
			c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
			return
		}

		// show metric
		c.String(http.StatusOK, strconv.FormatInt(counterVal, 10))
	default:
		c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
		return
	}

}

func (s *Handler) OutputAllMetrics(c *gin.Context) {
	data, err := s.srvUsecase.GetAllDataUsecase()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Metrics": data,
	})
}
