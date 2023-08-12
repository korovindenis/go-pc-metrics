package main

import (
	"log"
	"os"

	"github.com/korovindenis/go-pc-metrics/internal/adapter/config"
	"github.com/korovindenis/go-pc-metrics/internal/adapter/logger"
	agent "github.com/korovindenis/go-pc-metrics/internal/agent/agentapp"
	agentUsecase "github.com/korovindenis/go-pc-metrics/internal/domain/usecase/agent"
)

const (
	ExitSucces    = 0
	ExitWithError = 1
)

func main() {
	// init config (flags and env)
	cfg, err := config.New()
	if err != nil {
		log.Println(err)
		os.Exit(ExitWithError)
	}

	// init logger
	stdLog := log.New(log.Writer(), "", log.Flags())
	log := logger.New(stdLog)

	// init usecases
	agentUsecase, err := agentUsecase.New()
	if err != nil {
		log.Println(err)
		os.Exit(ExitWithError)
	}

	// run agent
	if err := agent.Exec(agentUsecase, log, cfg); err != nil {
		log.Println(err)
		os.Exit(ExitWithError)
	}
	os.Exit(ExitSucces)
}
