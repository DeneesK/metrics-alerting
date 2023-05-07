package main

import (
	"fmt"

	"github.com/DeneesK/metrics-alerting/internal/services/metriccollector"
	"github.com/DeneesK/metrics-alerting/internal/services/urlpreparer"
	"github.com/levigross/grequests"
)

var (
	counterMetric string = "counter"
	gaugeMetric   string = "gauge"
	contentType   string = "text/plain"
	pollCount     string = "PollCount"
	randomValue   string = "RandomValue"
)

type Collector interface {
	StartCollect()
	GetRuntimeMetrics() metriccollector.RuntimeMetrics
	GetRandomValue() float64
	GetPollCount() int64
}

func sendMetrics(ms Collector) error {
	runtimeMetrics := ms.GetRuntimeMetrics()
	cpuMetrics := runtimeMetrics.GetCPUMetrics()
	memMetrics := runtimeMetrics.GetMemMetrics()

	ro := grequests.RequestOptions{Headers: map[string]string{"Content-Type": contentType}}
	session := grequests.NewSession(&ro)
	defer session.CloseIdleConnections()

	for k, v := range cpuMetrics {
		if err := send(session, gaugeMetric, k, v); err != nil {
			return err
		}
	}
	for k, v := range memMetrics {
		if err := send(session, gaugeMetric, k, v); err != nil {
			return err
		}
	}
	if err := send(session, gaugeMetric, randomValue, ms.GetRandomValue()); err != nil {
		return err
	}

	if err := send(session, counterMetric, pollCount, ms.GetPollCount()); err != nil {
		return err
	}
	return nil
}

func send(session *grequests.Session, metricType string, metricName string, value interface{}) error {
	url, err := urlpreparer.PrepareURL(runAddr, metricType, metricName, value)
	if err != nil {
		return fmt.Errorf("unable to send report: %w", err)
	} else {
		_, err = postReport(session, url)
		if err != nil {
			return fmt.Errorf("unable to send report: %w", err)
		}
	}
	return nil
}

func postReport(s *grequests.Session, url string) (*grequests.Response, error) {
	return s.Post(url, s.RequestOptions)
}
