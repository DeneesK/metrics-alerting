package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	metricsStorage, err := storage.NewStorage(conf.filePath, conf.storeInterval, conf.isRestore, log, conf.dsn)
	if err != nil {
		return err
	}
	defer metricsStorage.Close()

	var wg sync.WaitGroup

	r := api.Routers(metricsStorage, log, []byte(conf.hashKey))
	srv := http.Server{Addr: conf.runAddr, Handler: r}

	wg.Add(1)
	go runServer(&wg, &srv)
	log.Infof("server started at %s", conf.runAddr)

	<-termChan

	err = srv.Shutdown(context.Background())
	if err != nil {
		return nil
	}

	wg.Wait()

	return nil
}

func runServer(wg *sync.WaitGroup, srv *http.Server) {
	defer wg.Done()
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
