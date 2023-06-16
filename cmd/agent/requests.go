package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/DeneesK/metrics-alerting/internal/models"
	"github.com/DeneesK/metrics-alerting/internal/services/metriccollector"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	retryMax             int           = 3
	retryWaitMin         time.Duration = time.Second * 1
	retryWaitMax         time.Duration = time.Second * 5
	counterMetric        string        = "counter"
	gaugeMetric          string        = "gauge"
	contentType          string        = "application/json"
	pollCount            string        = "PollCount"
	randomValue          string        = "RandomValue"
	encodeType           string        = "gzip"
	additionalMetricsLen int           = 2
)

type Collector interface {
	StartCollect()
	GetRuntimeMetrics() metriccollector.RuntimeMetrics
	GetRandomValue() float64
	GetPollCount() int64
	ResetPollCount()
}

func sendMetrics(ms Collector, runAddr string) error {
	url, err := url.JoinPath("http://", runAddr, "updates", "/")
	if err != nil {
		return fmt.Errorf("during attempt to create url error ocurred - %w", err)
	}
	runtimeMetrics := ms.GetRuntimeMetrics()
	cpuMetrics := runtimeMetrics.GetCPUMetrics()
	memMetrics := runtimeMetrics.GetMemMetrics()
	length := len(cpuMetrics) + len(memMetrics) + additionalMetricsLen
	metrics := make([]models.Metrics, 0, length)

	retryClient := retryablehttp.NewClient()

	retryClient.RetryMax = retryMax
	retryClient.RetryWaitMin = retryWaitMin
	retryClient.RetryWaitMax = retryWaitMax

	for k, v := range cpuMetrics {
		metrics = append(metrics, models.Metrics{ID: k, MType: gaugeMetric, Value: &v})
	}
	for k, v := range memMetrics {
		vFloat64 := float64(v)
		metrics = append(metrics, models.Metrics{ID: k, MType: gaugeMetric, Value: &vFloat64})
	}
	randomV := ms.GetRandomValue()
	pollC := ms.GetPollCount()
	metrics = append(metrics, models.Metrics{ID: randomValue, MType: gaugeMetric, Value: &randomV})
	metrics = append(metrics, models.Metrics{ID: pollCount, MType: counterMetric, Delta: &pollC})

	statusCode, err := sendBatch(retryClient, url, metrics)
	if err != nil {
		return fmt.Errorf("all attempts to establish the connection have been run out, during attempts to send data error ocurred - %w, ", err)
	}
	if statusCode == http.StatusOK {
		ms.ResetPollCount()
	}
	return nil
}

func sendBatch(retryClient *retryablehttp.Client, url string, metrics []models.Metrics) (int, error) {
	res, err := json.Marshal(&metrics)
	if err != nil {
		return 0, fmt.Errorf("serialization error - %w", err)
	}
	r, err := compress(res)
	if err != nil {
		return 0, fmt.Errorf("compressing error - %w", err)
	}
	req, err := retryablehttp.NewRequest("POST", url, r)
	if err != nil {
		return 0, fmt.Errorf("request error - %w", err)
	}
	req.Header.Add("Accept-Encoding", encodeType)
	req.Header.Add("Content-Encoding", encodeType)
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Content-Type", contentType)
	resp, err := retryClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
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
