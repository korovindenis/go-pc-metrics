package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	usecaseServer "github.com/korovindenis/go-pc-metrics/internal/domain/usecase/server"
)

type serverHandler struct {
	srvUsecase usecaseServer.IServerUsecase
}

type IServerHandler interface {
	ReceptionMetrics(c *gin.Context)
	OutputAllMetrics(c *gin.Context)
	OutputMetric(c *gin.Context)
}

func New(usecase usecaseServer.IServerUsecase) (IServerHandler, error) {
	return &serverHandler{
		srvUsecase: usecase,
	}, nil
}

func (s serverHandler) ReceptionMetrics(c *gin.Context) {
	// get metric prop from url
	namedURL := entity.ReqURI{
		MetricType: c.Param("metricType"),
		MetricName: c.Param("metricName"),
		MetricVal:  c.Param("metricVal"),
	}

	// run usecases
	switch namedURL.MetricType {
	case "gauge":
		// to float64
		metricVal, err := strconv.ParseFloat(namedURL.MetricVal, 64)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("metric value is wrong type"))
			return
		}

		// save metric
		if err = s.srvUsecase.SaveGauge(namedURL.MetricName, metricVal); err != nil {
			c.AbortWithError(http.StatusNotImplemented, errors.New("not Implemented server error"))
			return
		}
	case "counter":
		// to int64
		metricVal, err := strconv.ParseInt(namedURL.MetricVal, 10, 64)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("metric value is wrong type"))
			return
		}

		// save metric
		if err = s.srvUsecase.SaveCounter(namedURL.MetricName, metricVal); err != nil {
			c.AbortWithError(http.StatusNotImplemented, errors.New("not Implemented server error"))
			return
		}
	default:
		c.AbortWithError(http.StatusNotImplemented, errors.New("not Implemented server error"))
		return
	}

	c.Status(http.StatusOK)
}

func (s serverHandler) OutputMetric(c *gin.Context) {
	namedURL := entity.ReqURI{
		MetricType: c.Param("metricType"),
		MetricName: c.Param("metricName"),
	}
	switch namedURL.MetricType {
	case "gauge":
		// get metric
		gaugeVal, err := s.srvUsecase.GetGauge(namedURL.MetricName)
		if err != nil {
			if err == entity.ErrMetricNotFound {
				c.AbortWithError(http.StatusNotFound, errors.New("metric not found"))
				return
			}
			c.AbortWithError(http.StatusNotImplemented, errors.New("not Implemented server error"))
			return
		}

		// show metric
		c.String(http.StatusOK, strconv.FormatFloat(gaugeVal, 'f', 3, 64))
	case "counter":
		// get metric
		counterVal, err := s.srvUsecase.GetCounter(namedURL.MetricName)
		if err != nil {
			if err == entity.ErrMetricNotFound {
				c.AbortWithError(http.StatusNotFound, errors.New("metric not found"))
				return
			}
			c.AbortWithError(http.StatusNotImplemented, errors.New("not Implemented server error"))
			return
		}

		// show metric
		c.String(http.StatusOK, strconv.FormatInt(counterVal, 10))
	default:
		c.AbortWithError(http.StatusNotImplemented, errors.New("not Implemented server error"))
		return
	}

}

func (s serverHandler) OutputAllMetrics(c *gin.Context) {
	data, err := s.srvUsecase.GetAllData()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("internal server error"))
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Metrics": data,
	})
}
