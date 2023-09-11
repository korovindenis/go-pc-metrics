package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Define a mock for the IServerUsecase interface
type mockServerUsecase struct {
	mock.Mock
}

func (m *mockServerUsecase) SaveAllDataUsecase(ctx context.Context, metrics []entity.Metrics) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockServerUsecase) SaveGaugeUsecase(ctx context.Context, gaugeName string, gaugeValue float64) error {
	args := m.Called(gaugeName, gaugeValue)
	return args.Error(0)
}

func (m *mockServerUsecase) SaveCounterUsecase(ctx context.Context, counterName string, counterValue int64) error {
	args := m.Called(counterName, counterValue)
	return args.Error(0)
}

func (m *mockServerUsecase) GetGaugeUsecase(ctx context.Context, gaugeName string) (float64, error) {
	args := m.Called(gaugeName)
	return args.Get(0).(float64), args.Error(1)
}

func (m *mockServerUsecase) GetCounterUsecase(ctx context.Context, counterName string) (int64, error) {
	args := m.Called(counterName)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockServerUsecase) GetAllDataUsecase(ctx context.Context) (entity.MetricsType, error) {
	args := m.Called()
	return args.Get(0).(entity.MetricsType), args.Error(1)
}

func (m *mockServerUsecase) Ping(ctx context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockServerUsecase) SaveAllDataBatchUsecase(ctx context.Context, metrics []entity.Metrics) error {
	args := m.Called()
	return args.Error(0)
}

func TestReceptionMetric(t *testing.T) {
	mockUsecase := new(mockServerUsecase)
	handler, _ := New(mockUsecase)

	router := gin.Default()
	router.POST("/update/:metricType/:metricName/:metricVal", handler.ReceptionMetric)

	// t.Run("SaveGauge Success", func(t *testing.T) {
	// 	mockUsecase.On("SaveGaugeUsecase", "OtherSys", 471728.0).Return(nil).Once()

	// 	w := httptest.NewRecorder()
	// 	req, _ := http.NewRequest("POST", "/update/gauge/OtherSys/471728", nil)
	// 	router.ServeHTTP(w, req)

	// 	assert.Equal(t, http.StatusOK, w.Code)
	// 	mockUsecase.AssertExpectations(t)
	// })

	t.Run("SaveGauge Wrong Metric Value", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/update/gauge/OtherSys/wrongValue", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockUsecase.AssertNotCalled(t, "SaveGaugeUsecase")
	})

	mockUsecase.AssertExpectations(t)
}

func TestOutputMetric(t *testing.T) {
	mockUsecase := new(mockServerUsecase)
	handler, _ := New(mockUsecase)

	router := gin.Default()
	router.GET("/value/:metricType/:metricName", handler.OutputMetric)

	t.Run("GetGauge Success", func(t *testing.T) {
		mockUsecase.On("GetGaugeUsecase", "OtherSys").Return(471728.0, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/value/gauge/OtherSys", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "471728", w.Body.String())
		mockUsecase.AssertExpectations(t)
	})

	t.Run("GetGauge Metric Not Found", func(t *testing.T) {
		mockUsecase.On("GetGaugeUsecase", "InvalidMetric").Return(0.0, entity.ErrMetricNotFound).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/value/gauge/InvalidMetric", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockUsecase.AssertExpectations(t)
	})

	mockUsecase.AssertExpectations(t)
}

func TestOutputAllMetrics(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %s", err)
	}

	mockUsecase := new(mockServerUsecase)
	handler, _ := New(mockUsecase)
	router := gin.Default()
	router.LoadHTMLGlob(filepath.Dir(currentDir) + "/templates/*")

	router.GET("/", handler.OutputAllMetrics)

	t.Run("GetAllData Success", func(t *testing.T) {
		mockMetrics := entity.MetricsType{
			Gauge: map[string]float64{
				"Metric1": 123.45,
				"Metric2": 67.89,
			},
			Counter: map[string]int64{
				"Counter1": 100,
				"Counter2": 200,
			},
		}
		mockUsecase.On("GetAllDataUsecase").Return(mockMetrics, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Add assertions for response body as needed
		mockUsecase.AssertExpectations(t)
	})

	mockUsecase.AssertExpectations(t)
}
