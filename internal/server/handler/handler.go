package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
)

// function usecase
type usecase interface {
	SaveGaugeUsecase(ctx context.Context, gaugeName string, gaugeValue float64) error
	GetGaugeUsecase(ctx context.Context, gaugeName string) (float64, error)

	SaveCounterUsecase(ctx context.Context, counterName string, counterValue int64) error
	GetCounterUsecase(ctx context.Context, counterName string) (int64, error)

	SaveAllDataUsecase(ctx context.Context, metrics []entity.Metrics) error
	GetAllDataUsecase(ctx context.Context) (entity.MetricsType, error)

	SaveAllDataBatchUsecase(ctx context.Context, metrics []entity.Metrics) error

	Ping(ctx context.Context) error
}

// type usecasesWithBd interface {
// 	SaveAllDataBatchUsecase(ctx context.Context, metrics []entity.Metrics) error
// }

type Handler struct {
	serverUsecase usecase
}

func New(u usecase) (*Handler, error) {
	return &Handler{
		serverUsecase: u,
	}, nil
}

func (s *Handler) ReceptionMetric(c *gin.Context) {
	var metrics entity.Metrics
	ctx := c.Request.Context()

	if c.GetHeader("Content-Type") == "application/json" {
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
		if c.Param("metricType") == "counter" {
			metricVal, err := strconv.ParseInt(c.Param("metricVal"), 10, 64)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, entity.ErrInputVarIsWrongType)
				return
			}

			metrics.Delta = &metricVal
		}
		if c.Param("metricType") == "gauge" {
			metricVal, err := strconv.ParseFloat(c.Param("metricVal"), 64)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, entity.ErrInputVarIsWrongType)
				return
			}

			metrics.Value = &metricVal
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
		if err := s.serverUsecase.SaveGaugeUsecase(ctx, metrics.ID, *metrics.Value); err != nil {
			c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
			return
		}

		// show actual metrics
		gaugeVal, _ := s.serverUsecase.GetGaugeUsecase(ctx, metrics.ID)
		metrics.Value = &gaugeVal
		if c.Param("metricType") == "" {
			c.JSON(http.StatusOK, metrics)
			return
		}
		c.Status(http.StatusOK)
	case "counter":
		// save metric
		if err := s.serverUsecase.SaveCounterUsecase(ctx, metrics.ID, *metrics.Delta); err != nil {
			c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
			return
		}

		// show actual metrics
		counterVal, _ := s.serverUsecase.GetCounterUsecase(ctx, metrics.ID)
		metrics.Delta = &counterVal
		if c.Param("metricType") == "" {
			c.JSON(http.StatusOK, metrics)
			return
		}
		c.Status(http.StatusOK)
	default:
		c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
		return
	}

}
func (s *Handler) ReceptionMetrics(c *gin.Context) {
	var metrics []entity.Metrics
	ctx := c.Request.Context()

	if c.GetHeader("Content-Type") != "application/json" {
		c.JSON(http.StatusBadRequest, entity.ErrInvalidURLFormat)
		return
	}
	if err := c.ShouldBindJSON(&metrics); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, entity.ErrInvalidURLFormat)
		return
	}
	// usecasesWithBd, ok := s.serverUsecase.(usecasesWithBd)
	// if !ok {
	// 	c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
	// 	return
	// }

	// if err := usecasesWithBd.SaveAllDataBatchUsecase(ctx, metrics); err != nil {
	// 	c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
	// 	return
	// }

	if err := s.serverUsecase.SaveAllDataBatchUsecase(ctx, metrics); err != nil {
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}

	c.Status(200)
}

func (s *Handler) OutputMetric(c *gin.Context) {
	var metrics entity.Metrics
	ctx := c.Request.Context()

	if c.Param("metricType") == "" {
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
		fmt.Println("Handler  data" + fmt.Sprintf("%+v", metrics))

		// get metric
		gaugeVal, err := s.serverUsecase.GetGaugeUsecase(ctx, metrics.ID)
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
		if c.Param("metricType") == "" {
			c.JSON(http.StatusOK, metrics)
			return
		}
		c.String(http.StatusOK, strconv.FormatFloat(gaugeVal, 'g', -1, 64))
	case "counter":
		// get metric
		counterVal, err := s.serverUsecase.GetCounterUsecase(ctx, metrics.ID)
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
		if c.Param("metricType") == "" {
			c.JSON(http.StatusOK, metrics)
			return
		}
		c.String(http.StatusOK, strconv.FormatInt(counterVal, 10))
	default:
		c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
		return
	}
}

func (s *Handler) OutputAllMetrics(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := s.serverUsecase.GetAllDataUsecase(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Metrics": data,
	})
}

func (s *Handler) Ping(c *gin.Context) {
	if err := s.serverUsecase.Ping(c.Request.Context()); err != nil {
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
