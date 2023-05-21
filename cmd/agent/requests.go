package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/DeneesK/metrics-alerting/internal/models"
	"github.com/DeneesK/metrics-alerting/internal/services/metriccollector"
	"github.com/levigross/grequests"
)

var (
	counterMetric string  = "counter"
	gaugeMetric   string  = "gauge"
	contentType   string  = "application/json"
	pollCount     string  = "PollCount"
	randomValue   string  = "RandomValue"
	cvalue        int64   = 0
	gvalue        float64 = 0
)

type Collector interface {
	StartCollect()
	GetRuntimeMetrics() metriccollector.RuntimeMetrics
	GetRandomValue() float64
	GetPollCount() int64
	ResetPollCount()
}

func sendMetrics(ms Collector) error {
	url, err := url.JoinPath("http://", runAddr, "update", "/")
	if err != nil {
		return err
	}
	runtimeMetrics := ms.GetRuntimeMetrics()
	cpuMetrics := runtimeMetrics.GetCPUMetrics()
	memMetrics := runtimeMetrics.GetMemMetrics()

	ro := grequests.RequestOptions{Headers: map[string]string{"Content-Type": contentType}}
	session := grequests.NewSession(&ro)
	defer session.CloseIdleConnections()

	for k, v := range cpuMetrics {
		if _, err := send(session, url, gaugeMetric, k, v); err != nil {
			return err
		}
	}
	for k, v := range memMetrics {
		if _, err := send(session, url, gaugeMetric, k, v); err != nil {
			return err
		}
	}
	if _, err := send(session, url, gaugeMetric, randomValue, ms.GetRandomValue()); err != nil {
		return err
	}
	statusCode, err := send(session, url, counterMetric, pollCount, ms.GetPollCount())
	if err != nil {
		return err
	}
	if statusCode == http.StatusOK {
		ms.ResetPollCount()
	}
	return nil
}

func send(session *grequests.Session, url string, metricType string, metricName string, value interface{}) (int, error) {
	m := models.Metrics{Delta: &cvalue, Value: &gvalue}
	switch metricType {
	case "counter":
		switch t := value.(type) {
		case uint64:
			m.ID = metricName
			m.MType = metricType
			*m.Delta = int64(value.(uint64))
		case int64:
			m.ID = metricName
			m.MType = metricType
			*m.Delta = value.(int64)
		default:
			return 0, fmt.Errorf("unable to send report, value must be uint64, int64 or float64, have - %v", t)
		}
	case "gauge":
		switch t := value.(type) {
		case uint64:
			m.ID = metricName
			m.MType = metricType
			*m.Value = float64(value.(uint64))
		case float64:
			m.ID = metricName
			m.MType = metricType
			*m.Value = value.(float64)
		default:
			return 0, fmt.Errorf("unable to send report, value must be uint64 or float64, have - %v", t)
		}
	}
	res, err := json.Marshal(&m)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	resp, err := session.Post(url, &grequests.RequestOptions{JSON: res})
	if err != nil {
		return 0, fmt.Errorf("unable to send report: %w", err)
	}
	return resp.StatusCode, nil
}
