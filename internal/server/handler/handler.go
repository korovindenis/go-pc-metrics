package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	"github.com/korovindenis/go-pc-metrics/internal/encrypt"
)

//go:generate mockery --name usecase --exported
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

//go:generate mockery --name cfg --exported
type cfg interface {
	UseCryptoKey() bool
	GetKey() string
}

type Handler struct {
	serverUsecase usecase
	useCryptoKey  bool
	cryptoKey     string
}

func New(u usecase, cfg cfg) (*Handler, error) {
	return &Handler{
		serverUsecase: u,
		useCryptoKey:  cfg.UseCryptoKey(),
		cryptoKey:     cfg.GetKey(),
	}, nil
}

func (s *Handler) ReceptionMetric(c *gin.Context) {
	var metrics entity.Metrics
	ctx := c.Request.Context()

	if c.GetHeader("Content-Type") == "application/json" {
		// get metric from body
		if err := c.ShouldBindJSON(&metrics); err != nil {
			c.Error(fmt.Errorf("%s %w", "ReceptionMetric Json", err))
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
				c.Error(fmt.Errorf("%s %w", "ReceptionMetric ParseInt", err))
				c.AbortWithError(http.StatusBadRequest, entity.ErrInputVarIsWrongType)
				return
			}

			metrics.Delta = &metricVal
		}
		if c.Param("metricType") == "gauge" {
			metricVal, err := strconv.ParseFloat(c.Param("metricVal"), 64)
			if err != nil {
				c.Error(fmt.Errorf("%s %w", "ReceptionMetric ParseFloat", err))
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
			c.Error(fmt.Errorf("%s %w", "ReceptionMetric SaveGaugeUsecase", err))
			c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
			return
		}

		// show actual metrics
		gaugeVal, err := s.serverUsecase.GetGaugeUsecase(ctx, metrics.ID)
		if err != nil {
			c.Error(fmt.Errorf("%s %w", "ReceptionMetric GetGaugeUsecase", err))
			c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
			return
		}
		metrics.Value = &gaugeVal
		if c.Param("metricType") == "" {
			c.JSON(http.StatusOK, metrics)
			return
		}
		c.Status(http.StatusOK)
	case "counter":
		// save metric
		if err := s.serverUsecase.SaveCounterUsecase(ctx, metrics.ID, *metrics.Delta); err != nil {
			c.Error(fmt.Errorf("%s %w", "ReceptionMetric SaveCounterUsecase", err))
			c.AbortWithError(http.StatusNotImplemented, entity.ErrNotImplementedServerError)
			return
		}

		// show actual metrics
		counterVal, err := s.serverUsecase.GetCounterUsecase(ctx, metrics.ID)
		if err != nil {
			c.Error(fmt.Errorf("%s %w", "ReceptionMetric GetCounterUsecase", err))
			c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
			return
		}
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

	var requestBodyBuffer bytes.Buffer
	teeReader := io.TeeReader(c.Request.Body, &requestBodyBuffer)

	requestBody, _ := io.ReadAll(teeReader)
	defer c.Request.Body.Close()

	if s.useCryptoKey {
		decryptBody, err := encrypt.Decrypt(s.cryptoKey, string(requestBody))
		if err != nil {
			c.Error(fmt.Errorf("%s %w", "ReceptionMetrics DecryptData", err))
			c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
			return
		}
		requestBody = []byte(decryptBody)
	}

	if err := json.Unmarshal(requestBody, &metrics); err != nil {
		c.Error(fmt.Errorf("%s %w", "ReceptionMetrics Unmarshal", err))
		c.JSON(http.StatusBadRequest, entity.ErrInvalidURLFormat)
		return
	}

	if err := s.serverUsecase.SaveAllDataBatchUsecase(ctx, metrics); err != nil {
		c.Error(fmt.Errorf("%s %w", "ReceptionMetrics SaveAllDataBatchUsecase", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, struct{}{})
}

func (s *Handler) OutputMetric(c *gin.Context) {
	var metrics entity.Metrics
	ctx := c.Request.Context()

	if c.Param("metricType") == "" {
		// get metric from body
		if err := c.ShouldBindJSON(&metrics); err != nil {
			c.Error(fmt.Errorf("%s %w", "OutputMetric ShouldBindJSON", err))
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
		gaugeVal, err := s.serverUsecase.GetGaugeUsecase(ctx, metrics.ID)
		if err != nil {
			if errors.Is(err, entity.ErrMetricNotFound) {
				c.Error(fmt.Errorf("%s %w", "OutputMetric GetGaugeUsecase ErrMetricNotFound", err))
				c.AbortWithError(http.StatusNotFound, entity.ErrInputMetricNotFound)
				return
			}
			c.Error(fmt.Errorf("%s %w", "OutputMetric GetGaugeUsecase", err))
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
				c.Error(fmt.Errorf("%s %w", "OutputMetric GetCounterUsecase ErrInputMetricNotFound", err))
				c.AbortWithError(http.StatusNotFound, entity.ErrInputMetricNotFound)
				return
			}
			c.Error(fmt.Errorf("%s %w", "OutputMetric GetCounterUsecase", err))
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
		c.Error(fmt.Errorf("%s %w", "OutputAllMetrics GetAllDataUsecase", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Metrics": data,
	})
}

func (s *Handler) Ping(c *gin.Context) {
	if err := s.serverUsecase.Ping(c.Request.Context()); err != nil {
		c.Error(fmt.Errorf("%s %w", "Ping", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
