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

	// init bd
	storage, err := storage.New()
	if err != nil {
		log.Println(err)
		os.Exit(ExitWithError)
	}

	// init server usecases
	serverUsecase, err := serverusecase.New(storage)
	if err != nil {
		log.Println(err)
		os.Exit(ExitWithError)
	}

	// init server handlers
	serverHandler, err := serverhandler.New(serverUsecase)
	if err != nil {
		log.Println(err)
		os.Exit(ExitWithError)
	}

	// run web server
	if err := serverapp.Exec(cfg, serverHandler); err != nil {
		log.Println(err)
		os.Exit(ExitWithError)
	}
	os.Exit(ExitSucces)
}
