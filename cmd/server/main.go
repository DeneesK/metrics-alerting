package main

import (
	"log"
	"net/http"

	"github.com/DeneesK/metrics-alerting/cmd/server/memstorage"
	"github.com/go-chi/chi/v5"
)

func main() {
	parseFlags()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	metricsStorage := memstorage.NewMemStorage()
	r := Routers(metricsStorage)
	log.Printf("server started at %s", flagRunAddr)
	return http.ListenAndServe(flagRunAddr, r)
}

func Routers(ms memstorage.MemStorage) chi.Router {
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{value}", update(&ms))
	r.Get("/value/{metricType}/{metricName}", value(&ms))
	r.Get("/", metrics(&ms))
	return r
}
