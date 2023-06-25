package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/DeneesK/metrics-alerting/internal/services/metriccollector"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conf, err := parseFlags()

	if err != nil {
		return fmt.Errorf("args parse failed %w", err)
	}

	reportInterval := time.Duration(conf.reportingInterval) * time.Second
	ms := metriccollector.NewCollector(conf.pollingInterval)

	ch := make(chan metriccollector.RuntimeMetrics, conf.rateLimit)
	go ms.StartCollect()
	go ms.StartAdditionalCollect()
	go ms.FillChanel(ch)

	log.Printf("client started sending data on %s", conf.runAddr)
	var wg sync.WaitGroup
	for i := 0; i < conf.rateLimit; i++ {
		wg.Add(1)
		go func() {
			for metrics := range ch {
				if err := sendMetrics(metrics, &ms, conf.runAddr, conf.hashKey); err != nil {
					log.Println(err)
				}
				time.Sleep(reportInterval)
			}
		}()
	}
	wg.Wait()
	return nil
}
