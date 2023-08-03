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

type reqURI struct {
	metricType string
	metricName string
	metricVal  string
}

func (s serverHandler) ReceptionMetics(w http.ResponseWriter, r *http.Request) {
	// get metric prop from url
	origURL := strings.Split(r.RequestURI, "/")[2:]
	namedURL := reqURI{
		metricType: origURL[0],
		metricName: origURL[1],
		metricVal:  origURL[2],
	}

	// run usecases
	switch namedURL.metricType {
	case "gauge":
		metricVal, err := strconv.ParseFloat(namedURL.metricVal, 64)
		if err != nil {
			http.Error(w, "Metric value is wrong type!", http.StatusBadRequest)
			return
		}

		if err = s.srvUsecase.SaveGauge(namedURL.metricName, metricVal); err != nil {
			http.Error(w, "Not Implemented server error", http.StatusNotImplemented)
			return
		}
	case "counter":
		metricVal, err := strconv.ParseInt(namedURL.metricVal, 10, 64)
		if err != nil {
			http.Error(w, "Metric value is wrong type!", http.StatusBadRequest)
			return
		}

		if err = s.srvUsecase.SaveCounter(namedURL.metricName, metricVal); err != nil {
			http.Error(w, "Not Implemented server error", http.StatusNotImplemented)
			return
		}
	default:
		http.Error(w, "Not Implemented server error", http.StatusNotImplemented)
		return
	}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
