package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/korovindenis/go-pc-metrics/cmd/agent/app"
	"github.com/korovindenis/go-pc-metrics/internal/agent/config"
	agentUsecase "github.com/korovindenis/go-pc-metrics/internal/domain/usecases/agent"
	customLogger "github.com/korovindenis/go-pc-metrics/internal/logger"
	"go.uber.org/zap"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func init() {
	buildVersion = "N/A"
	buildDate = "N/A"
	buildDate = "N/A"
}

func main() {
	// for graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	// init ctx
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// init config (flags and env)
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config: %s\n", err)
	}

	// init logger
	logger, err := customLogger.New(cfg)
	if err != nil {
		log.Fatalf("logger: %s\n", err)
	}

	// init usecases
	agentUsecase, err := agentUsecase.New()
	if err != nil {
		logger.Fatal("init usecases", zap.Error(err))
	}

	// run agent
	go func() {
		if err := app.Run(ctx, agentUsecase, logger, cfg); err != nil {
			logger.Fatal("agent: ", zap.Error(err))
		}
	}()

	// graceful shutdown
	//wait for a signal to shutdown the server
	<-shutdown
	logger.Info("Shutting down...")

	// canceling the context to stop the app
	cancel()

	// we are waiting for the completion of the work of all goroutines
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	<-ctx.Done()

	logger.Info("Graceful shutdown complete")
}
