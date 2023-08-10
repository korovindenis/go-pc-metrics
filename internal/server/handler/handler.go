package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	usecaseServer "github.com/korovindenis/go-pc-metrics/internal/domain/usecase/server"
)

type ServerHandler interface {
	ReceptionMetrics(c *gin.Context)
	OutputAllMetrics(c *gin.Context)
	OutputMetric(c *gin.Context)
}

type Handler struct {
	srvUsecase usecaseServer.ServerUsecase
}

func New(usecase usecaseServer.ServerUsecase) (*Handler, error) {
	return &Handler{
		srvUsecase: usecase,
	}, nil
}

func (s Handler) ReceptionMetrics(c *gin.Context) {
	// get metric from url
	namedURL := entity.ReqURI{
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
		if err = s.srvUsecase.SaveGauge(namedURL.MetricName, metricVal); err != nil {
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
		if err = s.srvUsecase.SaveCounter(namedURL.MetricName, metricVal); err != nil {
			c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
			return
		}
	default:
		c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (s Handler) OutputMetric(c *gin.Context) {
	// get metric from url
	namedURL := entity.ReqURI{
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
		gaugeVal, err := s.srvUsecase.GetGauge(namedURL.MetricName)
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
		counterVal, err := s.srvUsecase.GetCounter(namedURL.MetricName)
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

func (s Handler) OutputAllMetrics(c *gin.Context) {
	data, err := s.srvUsecase.GetAllData()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Metrics": data,
	})
}
