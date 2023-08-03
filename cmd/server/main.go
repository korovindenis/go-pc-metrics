package main

import (
	"log"
	"os"

	storage "github.com/korovindenis/go-pc-metrics/internal/adapter/storage/memory"
	usecaseServer "github.com/korovindenis/go-pc-metrics/internal/domain/usecase/server"
	"github.com/korovindenis/go-pc-metrics/internal/http/handler"
	"github.com/korovindenis/go-pc-metrics/internal/http/server"
)

const httpAddress = "0.0.0.0:8080"

func main() {
	// init bd
	storage, err := storage.New()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// init server usecases
	srvUscs, err := usecaseServer.New(storage)
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
	if err := server.Exec(httpAddress, srvHdlrs); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
