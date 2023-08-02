package server

import (
	"net/http"

	"github.com/korovindenis/go-pc-metrics/internal/http/handler"
	"github.com/korovindenis/go-pc-metrics/internal/http/middleware"
)

// server main
func Exec(httpAddress string, handler handler.IServerHandler) error {
	mux := http.NewServeMux()
	// routes
	mux.Handle("/update/", middleware.CheckMethodAndContentType(http.HandlerFunc(handler.ReceptionMetics)))

	return http.ListenAndServe(httpAddress, mux)
}
