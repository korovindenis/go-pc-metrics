package main

import (
	"log"
	"os"

	"github.com/korovindenis/go-pc-metrics/cmd/agent/app"
	"github.com/korovindenis/go-pc-metrics/internal/agent/config"
	agentUsecase "github.com/korovindenis/go-pc-metrics/internal/domain/usecases/agent"
	customLogger "github.com/korovindenis/go-pc-metrics/internal/logger"
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
	logger, err := customLogger.New(cfg)
	if err != nil {
		log.Println(err)
		os.Exit(ExitWithError)
	}

	// init usecases
	agentUsecase, err := agentUsecase.New()
	if err != nil {
		logger.Error("init usecases", zap.Error(err))
		os.Exit(ExitWithError)
	}

	// run agent
	if err := app.Run(agentUsecase, logger, cfg); err != nil {
		logger.Error("run agent", zap.Error(err))
		os.Exit(ExitWithError)
	}
	os.Exit(ExitSucces)
}
