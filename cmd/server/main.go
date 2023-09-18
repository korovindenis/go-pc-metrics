package main

import (
	"context"
	"log"
	"os"

	"github.com/korovindenis/go-pc-metrics/cmd/server/app"
	"github.com/korovindenis/go-pc-metrics/internal/adapters/storage/disk"
	"github.com/korovindenis/go-pc-metrics/internal/adapters/storage/memory"
	database "github.com/korovindenis/go-pc-metrics/internal/adapters/storage/postgresql"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	serverusecase "github.com/korovindenis/go-pc-metrics/internal/domain/usecases/server"
	"github.com/korovindenis/go-pc-metrics/internal/logger"
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

func main() {
	// init config (flags and env)
	cfg, err := config.New()
	if err != nil {
		log.Println(err)
		os.Exit(ExitWithError)
	}

	// init logger
	err = logger.New(cfg)
	if err != nil {
		log.Println(err)
		os.Exit(ExitWithError)
	}

	// init storage
	var storage any
	switch cfg.GetStorageType() {
	case Bd:
		storage, err = database.New(cfg, logger.Log)
	case Disk:
		storage, err = disk.New(cfg, logger.Log)
	default:
		storage, err = memory.New(cfg, logger.Log)
	}

	if err != nil {
		logger.Log.Error("init storage", zap.Error(err))
		os.Exit(ExitWithError)
	}

	// init usecases
	serverUsecase, err := serverusecase.New(storage, cfg)
	if err != nil {
		logger.Log.Error("init usecases", zap.Error(err))
		os.Exit(ExitWithError)
	}

	// init handlers
	serverHandler, err := serverhandler.New(serverUsecase)
	if err != nil {
		logger.Log.Error("init handlers", zap.Error(err))
		os.Exit(ExitWithError)
	}

	// save to file
	if cfg.GetStorageType() == "disk" {
		ctx, cancel := context.WithCancel(context.Background())
		// Cancel the context when main() is terminated
		defer cancel()
		go serverUsecase.SaveAllDataUsecase(ctx, []entity.Metrics{})
	}

	// run web server
	if err := app.Run(cfg, serverHandler, logger.Log); err != nil {
		logger.Log.Error("run web server", zap.Error(err))
		os.Exit(ExitWithError)
	}
	os.Exit(ExitSucces)
}
