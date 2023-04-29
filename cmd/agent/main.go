package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	metric "github.com/DeneesK/metrics-alerting/cmd/agent/metriccollector"
	"github.com/DeneesK/metrics-alerting/cmd/agent/urlpreparer"
)

const (
	counterMetric  string        = "counter"
	gaugeMetric    string        = "gauge"
	reportInterval time.Duration = 10
	contentType    string        = "text/plain"
)

func sendReport(url string, contentType string) (resp *http.Response, err error) {
	buff := make([]byte, 0)
	return http.Post(url, contentType, bytes.NewBuffer(buff))
}

func sendMetrics(ms *metric.MetricStats) {
	time.Sleep(reportInterval * time.Second)
	metrics := urlpreparer.ParseNeededStats(ms.Stats)
	log.Println("sending... metric stats")
	for k, v := range metrics {
		url := urlpreparer.PrepareURL(k, gaugeMetric, v)
		sendReport(url, contentType)
	}
	sendReport(urlpreparer.PrepareURL("RandomValue", gaugeMetric, float32(ms.RandomValue)), contentType)
	sendReport(urlpreparer.PrepareURL("PollCount", counterMetric, float32(ms.PollCount)), contentType)
}

func main() {
	ms := metric.NewMetricStats()
	go ms.StartCollect()
	log.Println("metric collector started")
	for {
		sendMetrics(&ms)
	}
}
