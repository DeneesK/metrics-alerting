package main

import (
	"log"
	"net/http"

	"github.com/DeneesK/metrics-alerting/internal/api"
	"github.com/DeneesK/metrics-alerting/internal/logger"
	"github.com/DeneesK/metrics-alerting/internal/storage"
)

func main() {
	parseFlags()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	metricsStorage := storage.NewMemStorage()
	log, err := logger.LoggerInitializer(logLevel)
	if err != nil {
		return err
	}
	r := api.Routers(&metricsStorage)
	log.Infof("server started at %s", runAddr)
	return http.ListenAndServe(runAddr, r)
}
