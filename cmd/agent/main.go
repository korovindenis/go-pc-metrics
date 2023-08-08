package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/korovindenis/go-pc-metrics/internal/adapter/flags"
	"github.com/korovindenis/go-pc-metrics/internal/adapter/logger"
	agent "github.com/korovindenis/go-pc-metrics/internal/agent/agentapp"
	agentsecase "github.com/korovindenis/go-pc-metrics/internal/domain/usecase/agent"
)

func main() {
	// init flags
	_, err := flags.New()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// init flags for work tests :(
	var HTTPAddress string
	var reportInterval time.Duration
	var pollInterval time.Duration
	flag.StringVar(&HTTPAddress, "a", "localhost:8080", "HTTP server address")
	flag.DurationVar(&reportInterval, "r", 10*time.Second, "Report interval")
	flag.DurationVar(&pollInterval, "p", 2*time.Second, "Poll interval")
	flag.Parse()

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
	if err := agent.Exec(agntUscs, loggerInterface, "http://"+HTTPAddress, reportInterval, pollInterval); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
