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
	metricsStorage, err := storage.NewStorage(conf.filePath, conf.storeInterval, conf.isRestore, log, conf.dsn)
	if err != nil {
		return err
	}
	defer metricsStorage.Close()
	r := api.Routers(metricsStorage, log, conf.hashKey.Key)
	log.Infof("server started at %s", conf.runAddr)
	return http.ListenAndServe(conf.runAddr, r)
}
