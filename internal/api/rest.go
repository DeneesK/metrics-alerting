package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Store interface {
	Store(typeMetric, name, value string)
	GetValue(typeMetric, name string) string
	GetAll() string
}

func Routers(ms Store) chi.Router {
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{value}", update(ms))
	r.Get("/value/{metricType}/{metricName}", value(ms))
	r.Get("/", metrics(ms))
	return r
}

func update(storage Store) http.HandlerFunc {
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
		storage.Store(metricType, metricName, valueString)
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
	}
}

func value(storage Store) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		metricType := chi.URLParam(req, "metricType")
		metricName := chi.URLParam(req, "metricName")
		value := storage.GetValue(metricType, metricName)
		if value != "" {
			res.Write([]byte(value))
			res.Header().Add("Content-Type", "text/plain; charset=utf-8")
			res.WriteHeader(http.StatusOK)
			return
		}
		res.WriteHeader(http.StatusNotFound)
	}
}

func metrics(storage Store) http.HandlerFunc {
	return func(res http.ResponseWriter, _ *http.Request) {
		r := storage.GetAll()
		res.Write([]byte(r))
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
	}
}
