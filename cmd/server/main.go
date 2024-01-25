package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/korovindenis/go-pc-metrics/cmd/server/app"
	"github.com/korovindenis/go-pc-metrics/internal/adapters/storage/disk"
	"github.com/korovindenis/go-pc-metrics/internal/adapters/storage/memory"
	database "github.com/korovindenis/go-pc-metrics/internal/adapters/storage/postgresql"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	serverusecase "github.com/korovindenis/go-pc-metrics/internal/domain/usecases/server"
	customLogger "github.com/korovindenis/go-pc-metrics/internal/logger"
	"github.com/korovindenis/go-pc-metrics/internal/server/config"
	serverhandler "github.com/korovindenis/go-pc-metrics/internal/server/handler"
	"go.uber.org/zap"
)

const (
	Bd   = "database"
	Disk = "disk"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func init() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
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

	// init storage
	var storage any
	switch cfg.GetStorageType() {
	case Bd:
		storage, err = database.New(cfg, logger)
	case Disk:
		storage, err = disk.New(cfg, logger)
	default:
		storage, err = memory.New(cfg, logger)
	}

	if err != nil {
		logger.Fatal("init storage", zap.Error(err))
	}

	// init usecases
	serverUsecase, err := serverusecase.New(storage, cfg)
	if err != nil {
		logger.Fatal("init usecases", zap.Error(err))
	}

	// init handlers
	serverHandler, err := serverhandler.New(serverUsecase, cfg)
	if err != nil {
		logger.Fatal("init handlers", zap.Error(err))
	}

	// save to file
	if cfg.GetStorageType() == "disk" {
		ctx, cancel := context.WithCancel(context.Background())
		// Cancel the context when main() is terminated
		defer cancel()
		go serverUsecase.SaveAllDataUsecase(ctx, []entity.Metrics{})
	}

	go func() {
		// run web server
		if err := app.Run(ctx, cfg, serverHandler, logger); err != nil {
			logger.Error("run web server", zap.Error(err))
		}
	}()

	// graceful shutdown
	//wait for a signal to shutdown the server
	<-shutdown
	logger.Info("Shutting down...")

	// canceling the context to stop the app
	cancel()

	// we are waiting for the completion of the work of all goroutines
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	<-ctx.Done()

	logger.Info("Graceful shutdown complete")
}
