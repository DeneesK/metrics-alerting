package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	ms := metriccollector.NewCollector(conf.pollingInterval)

	ch := make(chan metriccollector.RuntimeMetrics, conf.rateLimit)
	reportInterval := time.Duration(conf.reportingInterval) * time.Second
	go ms.StartCollect()
	go ms.StartAdditionalCollect()
	go ms.FillChanel(ch, reportInterval)

	log.Printf("client started sending data on %s", conf.runAddr)
	var wg sync.WaitGroup
	for i := 0; i < conf.rateLimit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for metrics := range ch {
				if err := sendMetrics(metrics, &ms, conf.runAddr, conf.hashKey); err != nil {
					log.Println(err)
				}
			}
		}()
	}
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
	close(ch)
	log.Println("All workers were shuted down")
	wg.Wait()
	return nil
}
