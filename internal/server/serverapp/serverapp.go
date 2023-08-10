package serverapp

import (
	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-pc-metrics/internal/server/handler"
	"github.com/korovindenis/go-pc-metrics/internal/server/middleware"
)

// server main
func Exec(httpAddress string, handler handler.ServerHandler) error {
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
