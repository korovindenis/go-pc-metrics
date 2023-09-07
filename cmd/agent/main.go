package main

import (
	"log"
	"os"

	agent "github.com/korovindenis/go-pc-metrics/internal/agent/agentapp"
	"github.com/korovindenis/go-pc-metrics/internal/agent/config"
	agentUsecase "github.com/korovindenis/go-pc-metrics/internal/domain/usecases/agent"
	"github.com/korovindenis/go-pc-metrics/internal/logger"
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

	// init usecases
	agentUsecase, err := agentUsecase.New()
	if err != nil {
		logger.Log.Error("init usecases", zap.Error(err))
		os.Exit(ExitWithError)
	}

	// run agent
	if err := agent.Exec(agentUsecase, logger.Log, cfg); err != nil {
		logger.Log.Error("run agent", zap.Error(err))
		os.Exit(ExitWithError)
	}
	os.Exit(ExitSucces)
}
