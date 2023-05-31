package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/DeneesK/metrics-alerting/internal/models"
	"github.com/DeneesK/metrics-alerting/internal/services/metriccollector"
	"github.com/levigross/grequests"
)

const (
	counterMetric string = "counter"
	gaugeMetric   string = "gauge"
	contentType   string = "application/json"
	pollCount     string = "PollCount"
	randomValue   string = "RandomValue"
	encodeType    string = "gzip"
)

var (
	cvalue int64   = 0
	gvalue float64 = 0
)

type Collector interface {
	StartCollect()
	GetRuntimeMetrics() metriccollector.RuntimeMetrics
	GetRandomValue() float64
	GetPollCount() int64
	ResetPollCount()
}

func sendMetrics(ms Collector, runAddr string) error {
	url, err := url.JoinPath("http://", runAddr, "update", "/")
	if err != nil {
		return fmt.Errorf("during attempt to create url error ocurred - %v", err)
	}
	runtimeMetrics := ms.GetRuntimeMetrics()
	cpuMetrics := runtimeMetrics.GetCPUMetrics()
	memMetrics := runtimeMetrics.GetMemMetrics()

	ro := grequests.RequestOptions{Headers: map[string]string{
		"Accept-Encoding":  encodeType,
		"Content-Encoding": encodeType,
		"Content-Type":     contentType},
	}
	session := grequests.NewSession(&ro)
	defer session.CloseIdleConnections()

	for k, v := range cpuMetrics {
		if _, err := send(session, url, gaugeMetric, k, v); err != nil {
			return fmt.Errorf("during attempt to send data error ocurred - %v", err)
		}
	}
	for k, v := range memMetrics {
		if _, err := send(session, url, gaugeMetric, k, v); err != nil {
			return fmt.Errorf("during attempt to send data error ocurred - %v", err)
		}
	}
	if _, err := send(session, url, gaugeMetric, randomValue, ms.GetRandomValue()); err != nil {
		return err
	}
	statusCode, err := send(session, url, counterMetric, pollCount, ms.GetPollCount())
	if err != nil {
		return fmt.Errorf("during attempt to send data error ocurred - %v", err)
	}
	if statusCode == http.StatusOK {
		ms.ResetPollCount()
	}
	return nil
}

func send(session *grequests.Session, url string, metricType string, metricName string, value interface{}) (int, error) {
	m := models.Metrics{Delta: &cvalue, Value: &gvalue}
	m.ID = metricName
	m.MType = metricType
	switch metricType {
	case "counter":
		switch t := value.(type) {
		case uint64:
			*m.Delta = int64(value.(uint64))
		case int64:
			*m.Delta = value.(int64)
		default:
			return 0, fmt.Errorf("unable to send report, counter value must be uint64 or int64, have - %v", t)
		}
	case "gauge":
		switch t := value.(type) {
		case uint64:
			*m.Value = float64(value.(uint64))
		case float64:
			*m.Value = value.(float64)
		default:
			return 0, fmt.Errorf("unable to send report, gauge value must be uint64 or float64, have - %v", t)
		}
	default:
		return 0, fmt.Errorf("unable to send report, metricType must be counter or gauge, have - %v", metricType)
	}
	res, err := json.Marshal(&m)
	if err != nil {
		return 0, fmt.Errorf("serialisation error - %v", err)
	}
	r, err := compress(res)
	if err != nil {
		return 0, fmt.Errorf("compressing error - %v", err)
	}
	resp, err := session.Post(url, &grequests.RequestOptions{JSON: r})
	if err != nil {
		return 0, fmt.Errorf("unable to send report: %w", err)
	}
	defer resp.Close()
	return resp.StatusCode, nil
}

func compress(b []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		return nil, err
	}
	_, err = gz.Write(b)
	if err != nil {
		return nil, err
	}
	err = gz.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
