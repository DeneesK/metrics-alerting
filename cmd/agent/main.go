package main

import (
	"log"

	"github.com/DeneesK/metrics-alerting/internal/services/metriccollector"
)

func main() {
	parseFlags()
	ms := metriccollector.NewCollector(flagpolltInterval)
	go ms.StartCollect()
	log.Printf("client started sending data on %s", flagRunAddr)
	for {
		sendMetrics(&ms)
	}
}
