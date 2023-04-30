package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/DeneesK/metrics-alerting/cmd/server/memstorage"
	"github.com/go-chi/chi/v5"
)

type Repository interface {
	StoreMetrics(type_, name, value string)
	Value(type_, name string) string
	Metrics() string
}

func update(storage Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		metricType := chi.URLParam(req, "metricType")
		metricName := chi.URLParam(req, "metricName")
		valueString := chi.URLParam(req, "value")
		switch metricType {
		case "gauge":
			_, err := strconv.ParseFloat(valueString, 32)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		case "counter":
			_, err := strconv.Atoi(valueString)
			if err != nil {
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
	r := chi.NewRouter()
	metricsStorage := memstorage.NewMemStorage()
	r.Post("/update/{metricType}/{metricName}/{value}", update(&metricsStorage))
	r.Get("/value/{metricType}/{metricName}", value(&metricsStorage))
	r.Get("/", metrics(&metricsStorage))
	log.Println("server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func UpdateRouter(ms memstorage.MemStorage) chi.Router {
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{value}", update(&ms))
	return r
}
