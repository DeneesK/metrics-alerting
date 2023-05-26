package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DeneesK/metrics-alerting/internal/logger"
	"github.com/DeneesK/metrics-alerting/internal/models"
	"github.com/DeneesK/metrics-alerting/internal/storage"
	"github.com/go-chi/chi"
)

const (
	contentType     = "application/json"
	contentTypeText = "text/plain"
)

type Store interface {
	Store(typeMetric string, name string, value interface{}) error
	GetValue(typeMetric, name string) (storage.Result, bool, error)
	GetCounterMetrics() map[string]int64
	GetGaugeMetrics() map[string]float64
}

type apiLogger interface {
	Fatal(args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
}

var log apiLogger

func Routers(ms Store, logging apiLogger) chi.Router {
	log = logging
	r := chi.NewRouter()
	r.Use(logger.WithLogging)
	r.Use(gzipMiddleware)
	r.Post("/update/", UpdateJSON(ms))
	r.Post("/value/", ValueJSON(ms))
	r.Post("/update/{metricType}/{metricName}/{value}", Update(ms))
	r.Get("/value/{metricType}/{metricName}", Value(ms))
	r.Get("/", Metrics(ms))
	return r
}

func UpdateJSON(storage Store) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var metric models.Metrics
		if err := json.NewDecoder(req.Body).Decode(&metric); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		switch metric.MType {
		case "gauge":
			storage.Store(metric.MType, metric.ID, *metric.Value)
			resp, err := json.Marshal(&metric)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			res.Header().Add("Content-Type", contentType)
			res.WriteHeader(http.StatusOK)
			res.Write(resp)
		case "counter":
			storage.Store(metric.MType, metric.ID, *metric.Delta)
			value, ok, err := storage.GetValue(metric.MType, metric.ID)
			if ok {
				metric.Delta = &value.Counter
				resp, err := json.Marshal(&metric)
				if err != nil {
					res.WriteHeader(http.StatusBadRequest)
					return
				}
				res.Header().Add("Content-Type", contentType)
				res.WriteHeader(http.StatusOK)
				res.Write(resp)
			}
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		default:
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

func ValueJSON(storage Store) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var metric models.Metrics
		if err := json.NewDecoder(req.Body).Decode(&metric); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		value, ok, err := storage.GetValue(metric.MType, metric.ID)
		if err != nil {
			log.Error(err)
		} else if ok {
			switch metric.MType {
			case "counter":
				metric.Delta = &value.Counter
				resp, err := json.Marshal(&metric)
				if err != nil {
					res.WriteHeader(http.StatusBadRequest)
					return
				}
				res.Header().Add("Content-Type", contentType)
				res.WriteHeader(http.StatusOK)
				res.Write(resp)
			case "gauge":
				metric.Value = &value.Gauge
				resp, err := json.Marshal(&metric)
				if err != nil {
					res.WriteHeader(http.StatusBadRequest)
					return
				}
				res.Header().Add("Content-Type", contentType)
				res.WriteHeader(http.StatusOK)
				res.Write(resp)
			}
			return
		}
		res.WriteHeader(http.StatusNotFound)
	}
}

func Update(storage Store) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		metricType := chi.URLParam(req, "metricType")
		metricName := chi.URLParam(req, "metricName")
		valueString := chi.URLParam(req, "value")
		switch metricType {
		case "gauge":
			v, err := strconv.ParseFloat(valueString, 64)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			storage.Store(metricType, metricName, v)
		case "counter":
			v, err := strconv.ParseInt(valueString, 10, 64)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			storage.Store(metricType, metricName, v)
		default:
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Add("Content-Type", contentTypeText)
		res.WriteHeader(http.StatusOK)
	}
}

func Value(storage Store) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		metricType := chi.URLParam(req, "metricType")
		metricName := chi.URLParam(req, "metricName")
		value, ok, err := storage.GetValue(metricType, metricName)
		if ok {
			res.Header().Add("Content-Type", contentTypeText)
			res.WriteHeader(http.StatusOK)
			switch metricType {
			case "counter":
				res.Write([]byte(strconv.FormatInt(value.Counter, 10)))
			case "gauge":
				res.Write([]byte(strconv.FormatFloat(value.Gauge, byte(102), -1, 64)))
			}
			return
		}
		if err != nil {
			log.Error(err)
		}
		res.WriteHeader(http.StatusNotFound)
	}
}

func Metrics(storage Store) http.HandlerFunc {
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
		res.Header().Add("Content-Type", "text/html")
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(r))
	}
}
