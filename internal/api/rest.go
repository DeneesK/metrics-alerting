package api

import (
	"encoding/json"
	"fmt"
	"log"
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

func Routers(ms Store) chi.Router {
	r := chi.NewRouter()
	r.Use(logger.WithLogging)
	r.Post("/update/", updateJSON(ms))
	r.Post("/value/", valueJSON(ms))
	r.Post("/update/{metricType}/{metricName}/{value}", update(ms))
	r.Get("/value/{metricType}/{metricName}", value(ms))
	r.Get("/", metrics(ms))
	return r
}

func RouterWithoutLogger(ms Store) chi.Router {
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{value}", update(ms))
	r.Get("/value/{metricType}/{metricName}", value(ms))
	r.Post("/update/", updateJSON(ms))
	r.Post("/value/", valueJSON(ms))
	r.Get("/", metrics(ms))
	return r
}

func updateJSON(storage Store) http.HandlerFunc {
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
		res.WriteHeader(http.StatusOK)
	}
}

func valueJSON(storage Store) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var metric models.Metrics
		fmt.Println(metric.MType)
		fmt.Println(metric.ID)
		if err := json.NewDecoder(req.Body).Decode(&metric); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		value, ok, err := storage.GetValue(metric.MType, metric.ID)
		if ok {
			res.Header().Add("Content-Type", contentType)
			res.WriteHeader(http.StatusOK)
			switch metric.MType {
			case "counter":
				metric.Delta = &value.Counter
				resp, err := json.Marshal(&metric)
				if err != nil {
					res.WriteHeader(http.StatusBadRequest)
					return
				}
				res.Header().Add("Content-Type", contentType)
				res.Write(resp)
			case "gauge":
				metric.Value = &value.Gauge
				resp, err := json.Marshal(&metric)
				if err != nil {
					res.WriteHeader(http.StatusBadRequest)
					return
				}
				res.Header().Add("Content-Type", contentType)
				res.Write(resp)
			}
			return
		}
		if err != nil {
			log.Println(err)
		}
		res.WriteHeader(http.StatusNotFound)
	}
}

func update(storage Store) http.HandlerFunc {
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

func value(storage Store) http.HandlerFunc {
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
				res.Write([]byte(strconv.FormatFloat(value.Gauge, byte(102), -3, 64)))
			}
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
		res.Header().Add("Content-Type", contentType)
		res.Write([]byte(r))
		res.WriteHeader(http.StatusOK)
	}
}
