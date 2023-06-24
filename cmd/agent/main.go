package main

import (
	"log"
	"time"

	"github.com/DeneesK/metrics-alerting/internal/services/metriccollector"
)

func main() {
	conf, err := parseFlags()

	if err != nil {
		log.Fatal(err)
	}

	reportInterval := time.Duration(conf.reportingInterval) * time.Second
	ms := metriccollector.NewCollector(conf.pollingInterval)

	go ms.StartCollect()
	go ms.StartAdditionalCollect()

	log.Printf("client started sending data on %s", conf.runAddr)

	for {
		for i := 0; i < conf.rateLimit; i++ {
			go func() {
				if err := sendMetrics(&ms, conf.runAddr, conf.hashKey); err != nil {
					log.Println(err)
				}
			}()
		}
		time.Sleep(reportInterval)
	}
}
