package main

import (
	"log"
	"os"

	"github.com/korovindenis/go-pc-metrics/internal/adapter/flags"
	storage "github.com/korovindenis/go-pc-metrics/internal/adapter/storage/memory"
	serverusecase "github.com/korovindenis/go-pc-metrics/internal/domain/usecase/server"
	"github.com/korovindenis/go-pc-metrics/internal/server/handler"
	server "github.com/korovindenis/go-pc-metrics/internal/server/serverapp"
)

func main() {
	// init flags
	config, err := flags.New()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// init bd
	storage, err := storage.New()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// init server usecases
	srvUscs, err := serverusecase.New(storage)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// init server handlers
	srvHdlrs, err := handler.New(srvUscs)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// run web server
	if err := server.Exec(config.GetHttpAddress(), srvHdlrs); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
