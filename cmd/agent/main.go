package main

import (
	"log"
	"time"

	metric "github.com/DeneesK/metrics-alerting/cmd/agent/metriccollector"
	"github.com/DeneesK/metrics-alerting/cmd/agent/urlpreparer"
	"github.com/levigross/grequests"
)

const (
	counterMetric  string        = "counter"
	gaugeMetric    string        = "gauge"
	reportInterval time.Duration = 10
	contentType    string        = "text/plain"
)

func sendReport(s *grequests.Session, url string) (*grequests.Response, error) {
	return s.Post(url, s.RequestOptions)
}

func sendMetrics(ms *metric.MetricStats) {
	time.Sleep(reportInterval * time.Second)
	metrics := urlpreparer.ParseNeededStats(ms.Stats)
	ro := grequests.RequestOptions{Headers: map[string]string{"Content-Type": contentType}}
	session := grequests.NewSession(&ro)
	defer session.CloseIdleConnections()
	for k, v := range metrics {
		url := urlpreparer.PrepareURL(k, gaugeMetric, v)
		sendReport(session, url)
	}
	sendReport(session, urlpreparer.PrepareURL("RandomValue", gaugeMetric, float32(ms.RandomValue)))
	sendReport(session, urlpreparer.PrepareURL("PollCount", counterMetric, float32(ms.PollCount)))
}

func main() {
	ms := metric.NewMetricStats()
	go ms.StartCollect()
	log.Println("client started")
	for {
		sendMetrics(&ms)
	}
}
