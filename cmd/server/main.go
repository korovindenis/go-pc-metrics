package main

import (
	"log"
	"os"

	"github.com/korovindenis/go-pc-metrics/internal/adapters/config"
	storage "github.com/korovindenis/go-pc-metrics/internal/adapters/storage/disk"
	serverusecase "github.com/korovindenis/go-pc-metrics/internal/domain/usecase/server"
	"github.com/korovindenis/go-pc-metrics/internal/logger"
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
	cfg, err := config.New(true)
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

	// run web server
	if err := serverapp.Exec(cfg, serverHandler, storage, logger.Log); err != nil {
		logger.Log.Error("run web server", zap.Error(err))
		os.Exit(ExitWithError)
	}
	os.Exit(ExitSucces)
}
