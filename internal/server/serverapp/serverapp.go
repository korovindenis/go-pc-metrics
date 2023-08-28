package serverapp

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-pc-metrics/internal/logger"
	"github.com/korovindenis/go-pc-metrics/internal/server/middleware"
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
}

// server main
func Exec(cfg config, handler serverHandler) error {
	httpAddress := cfg.GetServerAddress()
	router := gin.Default()

	// html template
	router.LoadHTMLGlob("./internal/server/templates/*.html")

	// middleware
	router.Use(logger.RequestLogger())
	router.Use(middleware.CheckMethod())
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// routes
	router.GET("/", handler.OutputAllMetrics)
	router.GET("/value/:metricType/:metricName", handler.OutputMetric)
	router.POST("/value/", handler.OutputMetric)
	router.POST("/update/:metricType/:metricName/:metricVal", handler.ReceptionMetrics)
	router.POST("/update/", handler.ReceptionMetrics)

	// start server
	return router.Run(httpAddress)
}
