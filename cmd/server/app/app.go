package app

import (
	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-pc-metrics/internal/logger"
	"github.com/korovindenis/go-pc-metrics/internal/server/middleware"
	"go.uber.org/zap/zapcore"
)

// function handler
type serverHandler interface {
	ReceptionMetric(c *gin.Context)
	ReceptionMetrics(c *gin.Context)
	OutputMetric(c *gin.Context)
	OutputAllMetrics(c *gin.Context)
	Ping(c *gin.Context)
}

// config functions
type cfg interface {
	GetServerAddress() string
}

// logger functions
type log interface {
	Info(msg string, fields ...zapcore.Field)
}

// server main
func Run(cfg cfg, handler serverHandler, log log) error {
	httpAddress := cfg.GetServerAddress()
	router := gin.Default()

	// html template
	router.LoadHTMLGlob("./internal/server/templates/*.html")

	// middleware
	router.Use(logger.RequestLogger())
	router.Use(middleware.CheckMethod())
	router.Use(gin.Recovery())
	router.Use(middleware.GzipMiddleware())
	router.Use(middleware.GzipResponseMiddleware())
	//router.Use(gzip.Gzip(gzip.DefaultCompression))

	// routes
	router.GET("/", handler.OutputAllMetrics)
	router.GET("/ping/", handler.Ping)
	router.GET("/value/:metricType/:metricName", handler.OutputMetric)
	router.POST("/value/", handler.OutputMetric)
	router.POST("/update/:metricType/:metricName/:metricVal", handler.ReceptionMetric)
	router.POST("/update/", handler.ReceptionMetric)
	router.POST("/updates/", handler.ReceptionMetrics)

	// start server
	return router.Run(httpAddress)
}
