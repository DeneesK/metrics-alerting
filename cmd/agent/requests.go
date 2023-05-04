package main

import (
	"log"
	"time"

	"github.com/DeneesK/metrics-alerting/internal/services/metriccollector"
	"github.com/DeneesK/metrics-alerting/internal/services/urlpreparer"
	"github.com/levigross/grequests"
)

var (
	counterMetric  string        = "counter"
	gaugeMetric    string        = "gauge"
	reportInterval time.Duration = time.Second * time.Duration(ReportInterval)
	contentType    string        = "text/plain"
)

func sendReport(s *grequests.Session, url string) (*grequests.Response, error) {
	return s.Post(url, s.RequestOptions)
}

func sendMetrics(ms metriccollector.Collector) {
	time.Sleep(reportInterval * time.Second)
	metrics := ms.GetMetrics()
	ro := grequests.RequestOptions{Headers: map[string]string{"Content-Type": contentType}}
	session := grequests.NewSession(&ro)
	defer session.CloseIdleConnections()
	for k, v := range metrics {
		url, err := urlpreparer.PrepareURL(RunAddr, gaugeMetric, k, v)
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = sendReport(session, url)
		if err != nil {
			log.Println(err)
		}
	}
	url, err := urlpreparer.PrepareURL(RunAddr, "RandomValue", gaugeMetric, ms.GetRandomValue())
	if err != nil {
		log.Println(err)
	} else {
		_, err = sendReport(session, url)
		if err != nil {
			log.Println(err)
		}
	}
	url, err = urlpreparer.PrepareURL(RunAddr, "PollCount", counterMetric, float64(ms.GetPollCount()))
	if err != nil {
		log.Println(err)
	} else {
		_, err = sendReport(session, url)
		if err != nil {
			log.Println(err)
		}
	}
}
