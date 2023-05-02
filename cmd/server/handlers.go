package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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
	return func(res http.ResponseWriter, _ *http.Request) {
		r := storage.Metrics()
		res.Write([]byte(r))
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
	}
}
