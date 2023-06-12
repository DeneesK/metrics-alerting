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
	"github.com/levigross/grequests"
)

const (
	fstAttempt           time.Duration = time.Duration(1) * time.Second
	sndAttempt           time.Duration = fstAttempt * 3
	thirdAttempt         time.Duration = fstAttempt * 5
	counterMetric        string        = "counter"
	gaugeMetric          string        = "gauge"
	contentType          string        = "application/json"
	pollCount            string        = "PollCount"
	randomValue          string        = "RandomValue"
	encodeType           string        = "gzip"
	additionalMetricsLen int           = 2
)

var conAttempts = []time.Duration{fstAttempt, sndAttempt, thirdAttempt}

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
	ro := grequests.RequestOptions{Headers: map[string]string{
		"Accept-Encoding":  encodeType,
		"Content-Encoding": encodeType,
		"Content-Type":     contentType},
	}
	session := grequests.NewSession(&ro)
	defer session.CloseIdleConnections()
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

	statusCode, err := sendBanch(session, url, metrics)
	if err != nil {
		return fmt.Errorf("all attempts to establish a connection have been run out, during attempts to send data error ocurred - %v, ", err)
	}
	if statusCode == http.StatusOK {
		ms.ResetPollCount()
	}
	return nil
}

func sendBanch(session *grequests.Session, url string, metrics []models.Metrics) (int, error) {
	res, err := json.Marshal(&metrics)
	if err != nil {
		return 0, fmt.Errorf("serialization error - %w", err)
	}
	r, err := compress(res)
	if err != nil {
		return 0, fmt.Errorf("compressing error - %w", err)
	}
	resp, err := session.Post(url, &grequests.RequestOptions{JSON: r})
	if err != nil {
		for i, attempt := range conAttempts {
			time.Sleep(attempt)
			resp, err := session.Post(url, &grequests.RequestOptions{JSON: r})
			if err != nil && i < 2 {
				continue
			}
			defer resp.Close()
			if err != nil && i == 2 {
				return 0, fmt.Errorf("unable to send report: %w", err)
			}
		}
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
