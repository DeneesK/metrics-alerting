package main

import (
	"log"
	"net/http"

	"github.com/DeneesK/metrics-alerting/internal/api"
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
	r := api.Routers(&metricsStorage)
	log.Printf("server started at %s", flagRunAddr)
	return http.ListenAndServe(flagRunAddr, r)
}
