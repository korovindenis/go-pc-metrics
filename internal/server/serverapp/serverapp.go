package serverapp

import (
	"github.com/gin-gonic/gin"
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
	GetHTTPAddress() string
}

// server main
func Exec(cfg config, handler serverHandler) error {
	httpAddress := cfg.GetHTTPAddress()
	router := gin.Default()

	// html template
	router.LoadHTMLGlob("./internal/server/templates/*.html")

	// middleware
	router.Use(middleware.CheckMethod())
	router.Use(gin.Recovery())

	// routes
	router.GET("/", handler.OutputAllMetrics)
	router.GET("/value/:metricType/:metricName", handler.OutputMetric)
	router.POST("/update/:metricType/:metricName/:metricVal", handler.ReceptionMetrics)

	// start server
	return router.Run(httpAddress)
}
