package main

import (
	"log"
	"os"

	"github.com/korovindenis/go-pc-metrics/internal/adapter/env"
	"github.com/korovindenis/go-pc-metrics/internal/adapter/flags"
	"github.com/korovindenis/go-pc-metrics/internal/adapter/logger"
	agent "github.com/korovindenis/go-pc-metrics/internal/agent/agentapp"
	agentsecase "github.com/korovindenis/go-pc-metrics/internal/domain/usecase/agent"
)

func main() {
	// init env
	configEnv, err := env.New()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// init flags
	config, err := flags.New(configEnv)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// init logger
	stndrtLog := log.New(log.Writer(), "", log.Flags())
	loggerInterface := logger.New(stndrtLog)

	// init usecases
	agntUscs, err := agentsecase.New()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// run agent
	if err := agent.Exec(agntUscs, loggerInterface, config.GetHTTPAddressWithScheme(), config.GetPollInterval(), config.GetReportInterval()); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
