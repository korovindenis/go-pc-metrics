package logger

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger *zap.Logger
var once sync.Once

type cfg interface {
	GetLogsLevel() string
}

func New(config cfg) (*zap.Logger, error) {
	// for singletone
	once.Do(func() {
		lvl, err := zap.ParseAtomicLevel(config.GetLogsLevel())
		if err != nil {
			panic(err)
		}
		cfg := zap.NewProductionConfig()
		cfg.Level = lvl

		zl, err := cfg.Build()
		if err != nil {
			panic(err)
		}
		defer zl.Sync()

		logger = zl
	})

	return logger, nil
}

func RequestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		// Processing request
		ctx.Next()

		endTime := time.Now()

		logger.With(
			zap.Any("HTTP REQUEST", struct {
				METHOD  string
				URI     string
				STATUS  int
				LATENCY time.Duration
			}{
				ctx.Request.Method,
				ctx.Request.RequestURI,
				ctx.Writer.Status(),
				endTime.Sub(startTime),
			}),
		).Info("Request Logging")

		ctx.Next()
	}
}
