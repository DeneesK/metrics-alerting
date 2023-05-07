package main

import (
	"fmt"
	"log"
	"time"

	"github.com/DeneesK/metrics-alerting/internal/services/metriccollector"
	"github.com/DeneesK/metrics-alerting/internal/services/urlpreparer"
	"github.com/levigross/grequests"
)

var (
	counterMetric  string        = "counter"
	gaugeMetric    string        = "gauge"
	reportInterval time.Duration = time.Second * time.Duration(reportingInterval)
	contentType    string        = "text/plain"
)

type Collector interface {
	StartCollect()
	GetRuntimeMetrics() metriccollector.RuntimeMetrics
	GetRandomValue() float64
	GetPollCount() int64
}

func sendReport(s *grequests.Session, url string) (*grequests.Response, error) {
	return s.Post(url, s.RequestOptions)
}

func sendMetrics(ms Collector) error {
	runtimeMetrics := ms.GetRuntimeMetrics()
	cpuMetrics := runtimeMetrics.GetCPUMetrics()
	memMetrics := runtimeMetrics.GetMemMetrics()

	ro := grequests.RequestOptions{Headers: map[string]string{"Content-Type": contentType}}
	session := grequests.NewSession(&ro)
	defer session.CloseIdleConnections()

	for k, v := range cpuMetrics {
		url, err := urlpreparer.PrepareURL(runAddr, gaugeMetric, k, v)
		if err != nil {
			return fmt.Errorf("unable to send report: %w", err)
		}
		_, err = sendReport(session, url)
		if err != nil {
			return fmt.Errorf("unable to send report: %w", err)
		}
	}

	for k, v := range memMetrics {
		url, err := urlpreparer.PrepareURL(runAddr, gaugeMetric, k, v)
		if err != nil {
			return fmt.Errorf("unable to send report: %w", err)
		}
		_, err = sendReport(session, url)
		if err != nil {
			return fmt.Errorf("unable to send report: %w", err)
		}
	}

	url, err := urlpreparer.PrepareURL(runAddr, "RandomValue", gaugeMetric, ms.GetRandomValue())
	if err != nil {
		log.Println(err)
	} else {
		_, err = sendReport(session, url)
		if err != nil {
			return fmt.Errorf("unable to send report: %w", err)
		}
	}
	url, err = urlpreparer.PrepareURL(runAddr, "PollCount", counterMetric, ms.GetPollCount())
	if err != nil {
		log.Println(err)
	} else {
		_, err = sendReport(session, url)
		if err != nil {
			return fmt.Errorf("unable to send report: %w", err)
		}
	}
	return nil
}
