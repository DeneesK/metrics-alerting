package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/DeneesK/metrics-alerting/internal/api"
	"github.com/DeneesK/metrics-alerting/internal/logger"
	"github.com/DeneesK/metrics-alerting/internal/storage"
	"go.uber.org/zap"
)

func main() {
	if err := parseFlags(); err != nil {
		log.Fatal(err)
	}
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
	ok, err := strconv.ParseBool(isRestore)
	if err != nil {
		return err
	}
	if ok {
		if err := metricsStorage.LoadFromFile(filePath); err != nil {
			log.Infof("during attempt to load from file, error occurred: %v", err)
		}
	}
	if filePath != "" {
		go save(&metricsStorage, log)
	}
	return http.ListenAndServe(runAddr, r)
}

func save(m *storage.MemStorage, logger *zap.SugaredLogger) {
	storeInterval := time.Duration(storeInterval) * time.Second
	dir, _ := path.Split(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0666)
		if err != nil {
			logger.Error(err)
		}
	}
	for {
		time.Sleep(storeInterval)
		if err := m.SaveToFile(filePath); err != nil {
			logger.Error(err)
		}
	}
}
