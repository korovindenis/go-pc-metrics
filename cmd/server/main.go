package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	ExitSucces = iota
	ExitWithError
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
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	// init config (flags and env)
	cfg, err := config.New()
	if err != nil {
		log.Println(err)
		os.Exit(ExitWithError)
	}

	// init logger
	logger, err := customLogger.New(cfg)
	if err != nil {
		log.Println(err)
		os.Exit(ExitWithError)
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
		logger.Error("init storage", zap.Error(err))
		os.Exit(ExitWithError)
	}

	// init usecases
	serverUsecase, err := serverusecase.New(storage, cfg)
	if err != nil {
		logger.Error("init usecases", zap.Error(err))
		os.Exit(ExitWithError)
	}

	// init handlers
	serverHandler, err := serverhandler.New(serverUsecase, cfg)
	if err != nil {
		logger.Error("init handlers", zap.Error(err))
		os.Exit(ExitWithError)
	}

	// save to file
	if cfg.GetStorageType() == "disk" {
		ctx, cancel := context.WithCancel(context.Background())
		// Cancel the context when main() is terminated
		defer cancel()
		go serverUsecase.SaveAllDataUsecase(ctx, []entity.Metrics{})
	}

	appNew := app.New()
	if cfg.IsGrpc() {
		err = app.RunGrpc(cfg, appNew, serverHandler, logger)
	} else {
		err = app.RunHttp(cfg, appNew, serverHandler, logger)
	}
	if err != nil {
		logger.Error("run http server", zap.Error(err))
		os.Exit(ExitWithError)
	}

	// end work
	<-shutdown
	app.Stop(appNew)

	os.Exit(ExitSucces)
}
