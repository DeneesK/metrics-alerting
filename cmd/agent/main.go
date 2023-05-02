package main

import (
	"log"
	"time"

	metric "github.com/DeneesK/metrics-alerting/cmd/agent/metriccollector"
	"github.com/DeneesK/metrics-alerting/cmd/agent/urlpreparer"
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

func sendMetrics(ms *metric.MetricStats) {
	time.Sleep(reportInterval * time.Second)
	metrics := ParseNeededStats(ms.Stats)
	ro := grequests.RequestOptions{Headers: map[string]string{"Content-Type": contentType}}
	session := grequests.NewSession(&ro)
	defer session.CloseIdleConnections()
	for k, v := range metrics {
		url := urlpreparer.PrepareURL(flagRunAddr, gaugeMetric, k, v)
		sendReport(session, url)
	}
	sendReport(session, urlpreparer.PrepareURL(flagRunAddr, "RandomValue", gaugeMetric, float64(ms.RandomValue)))
	sendReport(session, urlpreparer.PrepareURL(flagRunAddr, "PollCount", counterMetric, float64(ms.PollCount)))
}

func main() {
	parseFlags()
	ms := metric.NewMetricStats(flagpolltInterval)
	go ms.StartCollect()
	log.Printf("client started sending data on %s", flagRunAddr)
	for {
		sendMetrics(&ms)
	}
}
