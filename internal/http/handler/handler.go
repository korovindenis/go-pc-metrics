package handler

import (
	"net/http"
	"strconv"
	"strings"

	usecaseServer "github.com/korovindenis/go-pc-metrics/internal/domain/usecase/server"
)

type serverHandler struct {
	srvUsecase usecaseServer.IServerUsecase
}

type IServerHandler interface {
	ReceptionMetics(w http.ResponseWriter, r *http.Request)
}

func New(usecase usecaseServer.IServerUsecase) (IServerHandler, error) {
	return &serverHandler{
		srvUsecase: usecase,
	}, nil
}

type reqUri struct {
	metricType string
	metricName string
	metricVal  string
}

func (s serverHandler) ReceptionMetics(w http.ResponseWriter, r *http.Request) {
	// get metric prop from url
	origUrl := strings.Split(r.RequestURI, "/")[2:]
	namedUrl := reqUri{
		metricType: origUrl[0],
		metricName: origUrl[1],
		metricVal:  origUrl[2],
	}

	// run usecases
	switch namedUrl.metricType {
	case "gauge":
		metricVal, err := strconv.ParseFloat(namedUrl.metricVal, 64)
		if err != nil {
			http.Error(w, "Metric value is wrong type!", http.StatusBadRequest)
			return
		}

		if err = s.srvUsecase.SaveGauge(namedUrl.metricName, metricVal); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	case "counter":
		metricVal, err := strconv.ParseInt(namedUrl.metricVal, 10, 64)
		if err != nil {
			http.Error(w, "Metric value is wrong type!", http.StatusBadRequest)
			return
		}

		if err = s.srvUsecase.SaveCounter(namedUrl.metricName, metricVal); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
