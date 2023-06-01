package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DeneesK/metrics-alerting/internal/models"
	"github.com/DeneesK/metrics-alerting/internal/storage"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
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

func Routers(ms Store, logging *zap.SugaredLogger) chi.Router {
	r := chi.NewRouter()
	r.Use(withLogging(logging))
	r.Use(gzipMiddleware(logging))
	r.Post("/update/", UpdateJSON(ms, logging))
	r.Post("/value/", ValueJSON(ms, logging))
	r.Post("/update/{metricType}/{metricName}/{value}", Update(ms, logging))
	r.Get("/value/{metricType}/{metricName}", Value(ms, logging))
	r.Get("/", Metrics(ms, logging))
	return r
}

func UpdateJSON(storage Store, log *zap.SugaredLogger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var metric models.Metrics
		if err := json.NewDecoder(req.Body).Decode(&metric); err != nil {
			log.Errorf("during attempt to deserializing error ocurred: %v", err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		switch metric.MType {
		case "gauge":
			storage.Store(metric.MType, metric.ID, *metric.Value)
			resp, err := json.Marshal(&metric)
			if err != nil {
				log.Errorf("during attempt to serializing error ocurred: %v", err)
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			res.Header().Add("Content-Type", contentType)
			res.WriteHeader(http.StatusOK)
			res.Write(resp)
		case "counter":
			storage.Store(metric.MType, metric.ID, *metric.Delta)
			value, ok, err := storage.GetValue(metric.MType, metric.ID)
			if err != nil || !ok {
				log.Errorf("unable to find value in storage, metric type exists: %v, ocurred error: %v", ok, err)
				res.WriteHeader(http.StatusNotFound)
				return
			}
			metric.Delta = &value.Counter
			resp, err := json.Marshal(&metric)
			if err != nil {
				log.Errorf("during attempt to serializing error ocurred: %v", err)
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			res.Header().Add("Content-Type", contentType)
			res.WriteHeader(http.StatusOK)
			res.Write(resp)
		default:
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

func ValueJSON(storage Store, log *zap.SugaredLogger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var metric models.Metrics
		if err := json.NewDecoder(req.Body).Decode(&metric); err != nil {
			log.Errorf("during body's decoding error ocurred: %v", err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		value, ok, err := storage.GetValue(metric.MType, metric.ID)
		if err != nil || !ok {
			log.Errorf("unable to find value in storage, metric type exists: %v, ocurred error: %v", ok, err)
			res.WriteHeader(http.StatusNotFound)
			return
		}
		switch metric.MType {
		case "counter":
			metric.Delta = &value.Counter
			resp, err := json.Marshal(&metric)
			if err != nil {
				log.Errorf("during attempt to serializing error ocurred: %v", err)
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			res.Header().Add("Content-Type", contentType)
			res.WriteHeader(http.StatusOK)
			res.Write(resp)
			return
		case "gauge":
			metric.Value = &value.Gauge
			resp, err := json.Marshal(&metric)
			if err != nil {
				log.Errorf("during attempt to serializing error ocurred: %v", err)
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			res.Header().Add("Content-Type", contentType)
			res.WriteHeader(http.StatusOK)
			res.Write(resp)
			return
		}
	}
}

func Update(storage Store, log *zap.SugaredLogger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		metricType := chi.URLParam(req, "metricType")
		metricName := chi.URLParam(req, "metricName")
		valueString := chi.URLParam(req, "value")
		switch metricType {
		case "gauge":
			v, err := strconv.ParseFloat(valueString, 64)
			if err != nil {
				log.Errorf("during attempt to parse value error ocurred: %v", err)
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			storage.Store(metricType, metricName, v)
		case "counter":
			v, err := strconv.ParseInt(valueString, 10, 64)
			if err != nil {
				log.Errorf("during attempt to parse value error ocurred: %v", err)
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

func Value(storage Store, log *zap.SugaredLogger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		metricType := chi.URLParam(req, "metricType")
		metricName := chi.URLParam(req, "metricName")
		value, ok, err := storage.GetValue(metricType, metricName)
		if err != nil || !ok {
			log.Errorf("unable to find value in storage, metric type exists: %v, ocurred error: %v", ok, err)
			res.WriteHeader(http.StatusNotFound)
			return
		}
		res.Header().Add("Content-Type", contentTypeText)
		res.WriteHeader(http.StatusOK)
		switch metricType {
		case "counter":
			res.Write([]byte(strconv.FormatInt(value.Counter, 10)))
		case "gauge":
			res.Write([]byte(strconv.FormatFloat(value.Gauge, 'f', -1, 64)))
		}
	}
}

func Metrics(storage Store, log *zap.SugaredLogger) http.HandlerFunc {
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
