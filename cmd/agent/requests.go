package main

import (
	"time"

	"github.com/DeneesK/metrics-alerting/internal/services/metriccollector"
	"github.com/DeneesK/metrics-alerting/internal/services/urlpreparer"
	"github.com/levigross/grequests"
)

var (
	counterMetric  string        = "counter"
	gaugeMetric    string        = "gauge"
	reportInterval time.Duration = time.Second * time.Duration(flagreportInterval)
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
		url := urlpreparer.PrepareURL(flagRunAddr, gaugeMetric, k, v)
		sendReport(session, url)
	}
	sendReport(session, urlpreparer.PrepareURL(flagRunAddr, "RandomValue", gaugeMetric, ms.GetRandomValue()))
	sendReport(session, urlpreparer.PrepareURL(flagRunAddr, "PollCount", counterMetric, float64(ms.GetPollCount())))
}
