package main

import (
	"log"
	"os"

	"github.com/korovindenis/go-pc-metrics/internal/adapter/config"
	storage "github.com/korovindenis/go-pc-metrics/internal/adapter/storage/memory"
	serverusecase "github.com/korovindenis/go-pc-metrics/internal/domain/usecase/server"
	serverhandler "github.com/korovindenis/go-pc-metrics/internal/server/handler"
	serverapp "github.com/korovindenis/go-pc-metrics/internal/server/serverapp"
)

const (
	EXIT_SUCCES     = 0
	EXIT_WITH_ERROR = 1
)

func main() {
	// init config (flags and env)
	cfg, err := config.New()
	if err != nil {
		log.Println(err)
		os.Exit(EXIT_WITH_ERROR)
	}

	// init bd
	storage, err := storage.New()
	if err != nil {
		log.Println(err)
		os.Exit(EXIT_WITH_ERROR)
	}

	// init server usecases
	serverUsecase, err := serverusecase.New(storage)
	if err != nil {
		log.Println(err)
		os.Exit(EXIT_WITH_ERROR)
	}

	// init server handlers
	serverHandlers, err := serverhandler.New(serverUsecase)
	if err != nil {
		log.Println(err)
		os.Exit(EXIT_WITH_ERROR)
	}

	// run web server
	if err := serverapp.Exec(cfg, serverHandlers); err != nil {
		log.Println(err)
		os.Exit(EXIT_WITH_ERROR)
	}
	os.Exit(EXIT_SUCCES)
}
