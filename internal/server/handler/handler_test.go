package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	"github.com/korovindenis/go-pc-metrics/internal/server/handler/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ReceptionMetric(t *testing.T) {
	usecase := mocks.NewUsecase(t)
	cfg := mocks.NewCfg(t)
	cfg.On("UseCryptoKey").Return(false)
	cfg.On("GetKey").Return("")
	handler, _ := New(usecase, cfg)
	router := gin.Default()
	router.GET("/update/:metricType/:metricName/:metricVal", handler.ReceptionMetric)

	tests := []struct {
		name           string
		url            string
		header         http.Header
		statusCode     int
		err            error
		getGaugeErr    error
		saveGaugeErr   error
		getCounterErr  error
		saveCounterErr error
	}{
		{
			name:       "positive getGaugeUsecase",
			url:        "/update/gauge/OtherSys/471728",
			statusCode: http.StatusOK,
		},
		{
			name:        "negative getGaugeUsecase",
			url:         "/update/gauge/OtherSys/471728",
			statusCode:  http.StatusInternalServerError,
			getGaugeErr: errors.New("err"),
		},
		{
			name:       "positive saveGaugeUsecase",
			url:        "/update/gauge/OtherSys/471728",
			statusCode: http.StatusOK,
		},
		{
			name:         "negative saveGaugeUsecase",
			url:          "/update/gauge/OtherSys/471728",
			statusCode:   http.StatusNotImplemented,
			saveGaugeErr: errors.New("err"),
		},
		{
			name:       "positive getCounterUsecase",
			url:        "/update/counter/OtherSys/471728",
			statusCode: http.StatusOK,
		},
		{
			name:          "negative getCounterUsecase",
			url:           "/update/counter/OtherSys/471728",
			statusCode:    http.StatusInternalServerError,
			getCounterErr: errors.New("err"),
		},
		{
			name:       "positive saveCounterUsecase",
			url:        "/update/counter/OtherSys/471728",
			statusCode: http.StatusOK,
		},
		{
			name:           "negative saveCounterUsecase",
			url:            "/update/counter/OtherSys/471728",
			statusCode:     http.StatusNotImplemented,
			saveCounterErr: errors.New("err"),
		},

		{
			name:       "negative",
			url:        "/update/gauge/",
			statusCode: http.StatusMovedPermanently,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			w := httptest.NewRecorder()

			saveGaugeUsecase := usecase.On("SaveGaugeUsecase", mock.Anything, mock.Anything, mock.Anything).Return(tt.saveGaugeErr).Maybe()
			getGaugeUsecase := usecase.On("GetGaugeUsecase", mock.Anything, mock.Anything).Return(float64(0), tt.getGaugeErr).Maybe()
			saveCounterUsecase := usecase.On("SaveCounterUsecase", mock.Anything, mock.Anything, mock.Anything).Return(tt.saveCounterErr).Maybe()
			getCounterUsecase := usecase.On("GetCounterUsecase", mock.Anything, mock.Anything).Return(int64(0), tt.getCounterErr).Maybe()

			// Act
			req, err := http.NewRequest(http.MethodGet, tt.url, http.NoBody)
			if err != nil {
				t.Fatal(err)
			}
			req.Header = tt.header
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)

			// Unset
			saveGaugeUsecase.Unset()
			getGaugeUsecase.Unset()
			saveCounterUsecase.Unset()
			getCounterUsecase.Unset()
		})
	}
}

func TestHandler_Ping(t *testing.T) {
	usecase := mocks.NewUsecase(t)
	cfg := mocks.NewCfg(t)
	cfg.On("UseCryptoKey").Return(false)
	cfg.On("GetKey").Return("")
	handler, _ := New(usecase, cfg)
	router := gin.Default()
	router.GET("/ping", handler.Ping)

	tests := []struct {
		name       string
		header     http.Header
		statusCode int
		err        error
	}{
		{
			name:       "positive",
			statusCode: http.StatusOK,
		},
		{
			name:       "negative ",
			statusCode: http.StatusInternalServerError,
			err:        errors.New("err"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			w := httptest.NewRecorder()

			ping := usecase.On("Ping", mock.Anything).Return(tt.err).Maybe()

			// Act
			req, err := http.NewRequest(http.MethodGet, "/ping", http.NoBody)
			if err != nil {
				t.Fatal(err)
			}
			req.Header = tt.header
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)

			// Unset
			ping.Unset()
		})
	}
}

func TestHandler_OutputAllMetrics(t *testing.T) {
	usecase := mocks.NewUsecase(t)
	cfg := mocks.NewCfg(t)
	cfg.On("UseCryptoKey").Return(false)
	cfg.On("GetKey").Return("")
	handler, _ := New(usecase, cfg)
	router := gin.Default()
	router.GET("/", handler.OutputAllMetrics)

	tests := []struct {
		name       string
		statusCode int
		err        error
	}{
		{
			name:       "negative",
			statusCode: http.StatusInternalServerError,
			err:        errors.New("err"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			w := httptest.NewRecorder()

			getAllDataUsecase := usecase.On("GetAllDataUsecase", mock.Anything).Return(entity.MetricsType{Gauge: make(entity.GaugeType), Counter: make(entity.CounterType)}, tt.err).Maybe()

			// Act
			req, err := http.NewRequest(http.MethodGet, "/", http.NoBody)
			if err != nil {
				t.Fatal(err)
			}
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)

			// Unset
			getAllDataUsecase.Unset()
		})
	}
}

func TestNew(t *testing.T) {
	mockUsecase := mocks.NewUsecase(t)
	cfg := mocks.NewCfg(t)
	cfg.On("UseCryptoKey").Return(false)
	cfg.On("GetKey").Return("")
	tests := []struct {
		name    string
		u       usecase
		want    *Handler
		wantErr bool
	}{
		{
			name: "positive",
			u:    mockUsecase,
			want: &Handler{serverUsecase: mockUsecase},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.u, cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestHandler_ReceptionMetrics(t *testing.T) {
	usecase := mocks.NewUsecase(t)
	cfg := mocks.NewCfg(t)
	cfg.On("UseCryptoKey").Return(false)
	cfg.On("GetKey").Return("")

	handler, _ := New(usecase, cfg)
	router := gin.Default()
	router.POST("/updates", handler.ReceptionMetrics)

	tests := []struct {
		name       string
		url        string
		header     http.Header
		statusCode int
		err        error
		args       []entity.Metrics
	}{
		{
			name:       "check content type getGaugeUsecase",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "check content type getGaugeUsecase",
			statusCode: http.StatusOK,
			header: http.Header{
				"Content-Type": {"application/json"},
			},
			args: []entity.Metrics{entity.Metrics{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			args, _ := json.Marshal(tt.args)
			w := httptest.NewRecorder()

			saveAllDataBatchUsecase := usecase.On("SaveAllDataBatchUsecase", mock.Anything, mock.Anything, mock.Anything).Return(tt.err).Maybe()

			// Act
			req, err := http.NewRequest(http.MethodPost, "/updates", bytes.NewBuffer([]byte(args)))
			if err != nil {
				t.Fatal(err)
			}
			req.Header = tt.header
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)

			// Unset
			saveAllDataBatchUsecase.Unset()
		})
	}
}

func TestHandler_OutputMetric(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		s    *Handler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.OutputMetric(tt.args.c)
		})
	}
}
