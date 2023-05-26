package main

import (
	"log"
	"net/http"

	"github.com/DeneesK/metrics-alerting/internal/api"
	"github.com/DeneesK/metrics-alerting/internal/logger"
	"github.com/DeneesK/metrics-alerting/internal/storage"
	"go.uber.org/zap"
)

func main() {
	conf, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}
	log, err := logger.LoggerInitializer(conf.logLevel)
	if err != nil {
		log.Fatal(err)
	}
	if err := run(conf, log); err != nil {
		log.Fatal(err)
	}
}

func run(conf Conf, log *zap.SugaredLogger) error {
	metricsStorage := storage.NewMemStorage(conf.filePath, conf.storeInterval, conf.isRestore, log)
	r := api.Routers(metricsStorage, log)
	log.Infof("server started at %s", conf.runAddr)
	return http.ListenAndServe(conf.runAddr, r)
}
