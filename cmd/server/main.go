package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/DeneesK/metrics-alerting/cmd/server/memstorage"
	"github.com/go-chi/chi/v5"
)

type Repository interface {
	StoreMetrics(typeMetric, name, value string)
	Value(typeMetric, name string) string
	Metrics() string
}

func update(storage Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		metricType := chi.URLParam(req, "metricType")
		metricName := chi.URLParam(req, "metricName")
		valueString := chi.URLParam(req, "value")
		switch metricType {
		case "gauge":
			if _, err := strconv.ParseFloat(valueString, 64); err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		case "counter":
			if _, err := strconv.Atoi(valueString); err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		default:
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.StoreMetrics(metricType, metricName, valueString)
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
	}
}

func value(storage Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		metricType := chi.URLParam(req, "metricType")
		metricName := chi.URLParam(req, "metricName")
		value := storage.Value(metricType, metricName)
		if value != "" {
			res.Write([]byte(value))
			res.Header().Add("Content-Type", "text/plain; charset=utf-8")
			res.WriteHeader(http.StatusOK)
			return
		}
		res.WriteHeader(http.StatusNotFound)
	}
}

func metrics(storage Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		r := storage.Metrics()
		res.Write([]byte(r))
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
	}
}

func main() {
	parseFlags()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func Routers(ms memstorage.MemStorage) chi.Router {
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{value}", update(&ms))
	r.Get("/value/{metricType}/{metricName}", value(&ms))
	r.Get("/", metrics(&ms))
	return r
}

func run() error {
	metricsStorage := memstorage.NewMemStorage()
	r := Routers(metricsStorage)
	log.Printf("server started at %s", flagRunAddr)
	return http.ListenAndServe(flagRunAddr, r)
}
