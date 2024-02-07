// The main application, the backend
package app

import (
	"net"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-pc-metrics/internal/logger"
	"github.com/korovindenis/go-pc-metrics/internal/server/middleware"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
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
	GetTrustedSubnet() string
	GetKey() string
}

// logger functions
type log interface {
	Info(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
}

func New() chan int {
	return make(chan int)
}

// http server
func RunHttp(cfg cfg, resultCh chan int, handler serverHandler, log log) error {
	secretKey := cfg.GetKey()
	httpAddress := cfg.GetServerAddress()
	trustedSubnet := cfg.GetTrustedSubnet()
	router := gin.Default()

	// html template
	router.LoadHTMLGlob("./internal/server/templates/*.html")

	// middleware
	router.Use(logger.RequestLogger())
	router.Use(gin.Recovery())
	router.Use(middleware.CheckMethod())
	router.Use(middleware.CheckRealIP(trustedSubnet))
	router.Use(middleware.ErrorLogging(log))
	router.Use(middleware.Gzip())
	router.Use(middleware.GzipResponse())
	if secretKey != "" {
		const patternSign = `^/updates?/$`

		router.Use(middleware.SetSign(secretKey, patternSign))
		router.Use(middleware.CheckSign(log, secretKey, patternSign))
	}

	// routes
	router.GET("/", handler.OutputAllMetrics)
	router.GET("/ping/", handler.Ping)
	router.GET("/value/:metricType/:metricName", handler.OutputMetric)
	router.POST("/value/", handler.OutputMetric)
	router.POST("/update/:metricType/:metricName/:metricVal", handler.ReceptionMetric)
	router.POST("/update/", handler.ReceptionMetric)
	router.POST("/updates/", handler.ReceptionMetrics)

	// add pprof
	pprof.Register(router)

	// start server
	return router.Run(httpAddress)
}
func RunGrpc(cfg cfg, resultCh chan int, handler serverHandler, log log) error {
	address := cfg.GetServerAddress()

	// Determine the port for the server
	listen, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	// Create a gRPC server without a registered service
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(),
	)

	// Handle incoming gRPC requests
	if err := grpcServer.Serve(listen); err != nil {
		return err
	}

	return nil
}
func Stop(resultCh chan int) {
	close(resultCh)
}
