package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

const (
	contentTypeText = "text/plain; charset=utf-8"
)

type Store interface {
	Store(typeMetric, name, value string) error
	GetValue(typeMetric, name string) (string, bool, error)
	GetCounterMetrics() map[string]int
	GetGaugeMetrics() map[string]float64
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
			if _, err := strconv.ParseInt(valueString, 10, 64); err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		default:
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.Store(metricType, metricName, valueString)
		res.Header().Add("Content-Type", contentTypeText)
		res.WriteHeader(http.StatusOK)
	}
}

func value(storage Store) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		metricType := chi.URLParam(req, "metricType")
		metricName := chi.URLParam(req, "metricName")
		value, ok, err := storage.GetValue(metricType, metricName)
		if ok {
			res.Write([]byte(value))
			res.Header().Add("Content-Type", contentTypeText)
			res.WriteHeader(http.StatusOK)
			return
		}
		if err != nil {
			log.Panicln(err)
		}
		res.WriteHeader(http.StatusNotFound)
	}
}

func metrics(storage Store) http.HandlerFunc {
	return func(res http.ResponseWriter, _ *http.Request) {
		c := storage.GetCounterMetrics()
		g := storage.GetGaugeMetrics()
		r := ""
		for k, v := range c {
			r += fmt.Sprintf("[%s]: %d\n", k, v)
		}
		for k, v := range g {
			r += fmt.Sprintf("[%s]: %g\n", k, v)
		}
		res.Write([]byte(r))
		res.Header().Add("Content-Type", contentTypeText)
		res.WriteHeader(http.StatusOK)
	}
}
