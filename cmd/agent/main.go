package main

import (
	"log"
	"time"

	"github.com/DeneesK/metrics-alerting/internal/services/metriccollector"
)

func main() {
	if err := parseFlags(); err != nil {
		log.Fatal(err)
	}

	reportInterval := time.Duration(reportingInterval) * time.Second
	ms := metriccollector.NewCollector(pollingInterval)
	go ms.StartCollect()
	log.Printf("client started sending data on %s", runAddr)

	for {
		if err := sendMetrics(&ms); err != nil {
			log.Println(err)
		}
		time.Sleep(reportInterval)
	}
}
