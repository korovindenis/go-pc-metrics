package handler

import (
	"errors"
	"net/http"

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
	serverUsecase usecase
}

func New(u usecase) (*Handler, error) {
	return &Handler{
		serverUsecase: u,
	}, nil
}

func (s *Handler) ReceptionMetrics(c *gin.Context) {
	var metrics entity.Metrics

	contentType := c.GetHeader("Content-Type")
	if contentType == "application/json" {
		// get metric from body
		if err := c.ShouldBindJSON(&metrics); err != nil {
			c.JSON(http.StatusBadRequest, entity.ErrInvalidURLFormat)
			return
		}
	} else {
		// get metric from url
		metrics = entity.Metrics{
			MType: c.Param("metricType"),
			ID:    c.Param("metricName"),
		}
	}

	// validate metrics
	if metrics.ID == "" || metrics.MType == "" {
		c.AbortWithError(http.StatusNotFound, entity.ErrInvalidURLFormat)
		return
	}

	// run usecases
	switch metrics.MType {
	case "gauge":
		// save metric
		if err := s.serverUsecase.SaveGaugeUsecase(metrics.ID, *metrics.Value); err != nil {
			c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
			return
		}

		// show actual metrics
		gaugeVal, _ := s.serverUsecase.GetGaugeUsecase(metrics.ID)
		metrics.Value = &gaugeVal
		c.JSON(http.StatusOK, metrics)
	case "counter":
		// save metric
		if err := s.serverUsecase.SaveCounterUsecase(metrics.ID, *metrics.Delta); err != nil {
			c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
			return
		}

		// show actual metrics
		counterVal, _ := s.serverUsecase.GetCounterUsecase(metrics.ID)
		metrics.Delta = &counterVal
		c.JSON(http.StatusOK, metrics)
	default:
		c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
		return
	}
}

func (s *Handler) OutputMetric(c *gin.Context) {
	var metrics entity.Metrics

	contentType := c.GetHeader("Content-Type")
	if contentType == "application/json" {
		// get metric from body
		if err := c.ShouldBindJSON(&metrics); err != nil {
			c.JSON(http.StatusBadRequest, entity.ErrInvalidURLFormat)
			return
		}
	} else {
		// get metric from url
		metrics = entity.Metrics{
			MType: c.Param("metricType"),
			ID:    c.Param("metricName"),
		}
	}

	// validate metrics
	if metrics.MType == "" || metrics.ID == "" {
		c.AbortWithError(http.StatusNotFound, entity.ErrInvalidURLFormat)
		return
	}

	switch metrics.MType {
	case "gauge":
		// get metric
		gaugeVal, err := s.serverUsecase.GetGaugeUsecase(metrics.ID)
		if err != nil {
			if errors.Is(err, entity.ErrMetricNotFound) {
				c.AbortWithError(http.StatusNotFound, entity.ErrInputMetricNotFound)
				return
			}
			c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
			return
		}

		// show metric
		metrics.Value = &gaugeVal
		c.JSON(http.StatusOK, metrics)
	case "counter":
		// get metric
		counterVal, err := s.serverUsecase.GetCounterUsecase(metrics.ID)
		if err != nil {
			if errors.Is(err, entity.ErrMetricNotFound) {
				c.AbortWithError(http.StatusNotFound, entity.ErrInputMetricNotFound)
				return
			}
			c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
			return
		}

		// show metric
		metrics.Delta = &counterVal
		c.JSON(http.StatusOK, metrics)
	default:
		c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
		return
	}

}

func (s *Handler) OutputAllMetrics(c *gin.Context) {
	data, err := s.serverUsecase.GetAllDataUsecase()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Metrics": data,
	})
}
