package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// for singletone
var Log *zap.Logger = zap.NewNop()

type cfg interface {
	GetLogsLevel() string
}

func New(config cfg) error {
	lvl, err := zap.ParseAtomicLevel(config.GetLogsLevel())
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	defer zl.Sync()

	Log = zl

	return nil
}

func RequestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		// Processing request
		ctx.Next()

		endTime := time.Now()

		Log.With(
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
