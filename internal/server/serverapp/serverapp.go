package serverapp

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-pc-metrics/internal/logger"
	"github.com/korovindenis/go-pc-metrics/internal/server/middleware"
	"go.uber.org/zap/zapcore"
)

// function handler
type serverHandler interface {
	ReceptionMetrics(c *gin.Context)
	OutputMetric(c *gin.Context)
	OutputAllMetrics(c *gin.Context)
}

// config functions
type config interface {
	GetServerAddress() string
	GetStoreInterval() time.Duration
}

// storege functions
type storage interface {
	SaveAllData() error
}

// logger functions
type log interface {
	Info(msg string, fields ...zapcore.Field)
}

// server main
func Exec(cfg config, handler serverHandler, storage storage, log log) error {
	httpAddress := cfg.GetServerAddress()
	router := gin.Default()

	// html template
	router.LoadHTMLGlob("./internal/server/templates/*.html")

	// middleware
	router.Use(logger.RequestLogger())
	router.Use(middleware.CheckMethod())
	//router.Use(gin.Recovery())
	//router.Use(gzip.Gzip(gzip.DefaultCompression))

	// routes
	//router.GET("/", handler.OutputAllMetrics)
	//router.GET("/value/:metricType/:metricName", handler.OutputMetric)
	router.POST("/value/", handler.OutputMetric)
	//router.POST("/update/:metricType/:metricName/:metricVal", handler.ReceptionMetrics)
	router.POST("/update/", handler.ReceptionMetrics)

	// save data to disk
	go saveAllData(cfg, storage, log)

	// start server
	return router.Run(httpAddress)
}

func saveAllData(cfg config, storage storage, log log) {
	sendTicker := time.NewTicker(cfg.GetStoreInterval())
	defer sendTicker.Stop()

	for range sendTicker.C {
		log.Info("save to file")
		storage.SaveAllData()
	}
}
