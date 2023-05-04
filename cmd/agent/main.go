package main

import (
	"log"

	"github.com/DeneesK/metrics-alerting/internal/services/metriccollector"
)

func main() {
	parseFlags()
	ms := metriccollector.NewCollector(PollInterval)
	go ms.StartCollect()
	log.Printf("client started sending data on %s", RunAddr)
	for {
		sendMetrics(&ms)
	}
}
