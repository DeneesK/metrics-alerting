package main

import (
	"log"

	metric "github.com/DeneesK/metrics-alerting/cmd/agent/metriccollector"
)

func main() {
	parseFlags()
	ms := metric.NewMetricStats(flagpolltInterval)
	go ms.StartCollect()
	log.Printf("client started sending data on %s", flagRunAddr)
	for {
		sendMetrics(&ms)
	}
}
