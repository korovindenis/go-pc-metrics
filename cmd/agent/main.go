package main

import (
	"log"
	"os"
	"time"

	"github.com/korovindenis/go-pc-metrics/internal/adapter/logger"
	agent "github.com/korovindenis/go-pc-metrics/internal/agent/agentapp"
	agentsecase "github.com/korovindenis/go-pc-metrics/internal/domain/usecase/agent"
)

const (
	pollInterval   = time.Second * 2
	reportInterval = time.Second * 10
	httpAddress    = "http://localhost:8080"
)

func main() {
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
	if err := agent.Exec(agntUscs, loggerInterface, httpAddress, pollInterval, reportInterval); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
