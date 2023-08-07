package serverapp

import (
	"net/http"

	"github.com/korovindenis/go-pc-metrics/internal/server/handler"
	"github.com/korovindenis/go-pc-metrics/internal/server/middleware"
)

// server main
func Exec(httpAddress string, handler handler.IServerHandler) error {
	mux := http.NewServeMux()
	// routes
	mux.Handle("/update/", middleware.CheckMethodAndContentType(http.HandlerFunc(handler.ReceptionMetics)))

	return http.ListenAndServe(httpAddress, mux)
}
