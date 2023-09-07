package main

import (
	"context"
	"log"
	"os"

	storage "github.com/korovindenis/go-pc-metrics/internal/adapters/storage/postgresql"
	serverusecase "github.com/korovindenis/go-pc-metrics/internal/domain/usecases/server"
	"github.com/korovindenis/go-pc-metrics/internal/logger"
	"github.com/korovindenis/go-pc-metrics/internal/server/config"
	serverhandler "github.com/korovindenis/go-pc-metrics/internal/server/handler"
	serverapp "github.com/korovindenis/go-pc-metrics/internal/server/serverapp"
	"go.uber.org/zap"
)

const (
	ExitSucces = iota
	ExitWithError
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

	// init bd
	storage, err := storage.New(cfg)
	if err != nil {
		logger.Log.Error("init bd", zap.Error(err))
		os.Exit(ExitWithError)
	}

	// init usecases
	serverUsecase, err := serverusecase.New(storage)
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

	ctx, cancel := context.WithCancel(context.Background())
	// Cancel the context when main() is terminated
	defer cancel()
	// save to file
	go serverUsecase.SaveAllDataUsecase(ctx, cfg)

	// run web server
	if err := serverapp.Exec(cfg, serverHandler, logger.Log); err != nil {
		logger.Log.Error("run web server", zap.Error(err))
		os.Exit(ExitWithError)
	}
	os.Exit(ExitSucces)
}
