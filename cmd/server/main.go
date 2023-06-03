package main

import (
	"log"
	"net/http"

	"github.com/DeneesK/metrics-alerting/internal/api"
	"github.com/DeneesK/metrics-alerting/internal/logger"
	"github.com/DeneesK/metrics-alerting/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conf, err := parseFlags()
	if err != nil {
		return err
	}
	log, err := logger.LoggerInitializer(conf.logLevel)
	if err != nil {
		return err
	}
	metricsStorage := storage.NewMemStorage(conf.filePath, conf.storeInterval, conf.isRestore, log, conf.postgresDSN)
	r := api.Routers(metricsStorage, log)
	log.Infof("server started at %s", conf.runAddr)
	return http.ListenAndServe(conf.runAddr, r)
}
