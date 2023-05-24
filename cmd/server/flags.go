package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var runAddr string
var logLevel string
var storeInterval int
var filePath string
var isRestore string

func parseFlags() error {
	flag.StringVar(&runAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&logLevel, "l", "info", "log level")
	flag.StringVar(&filePath, "f", "tmp/metrics-db.json", "path to store file")
	flag.StringVar(&isRestore, "r", "true", "load saved data")
	flag.IntVar(&storeInterval, "i", 300, "interval of storing data on disk")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		runAddr = envRunAddr
	}
	if envRunAddr := os.Getenv("LOG_LEVEL"); envRunAddr != "" {
		runAddr = envRunAddr
	}
	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		filePath = envFilePath
	}
	if envIsRestore := os.Getenv("RESTORE"); envIsRestore != "" {
		isRestore = envIsRestore
	}
	if envStoreInterval, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		fsi, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			err = fmt.Errorf("value of STORE_INTERVAL is not a integer: %w", err)
			return err
		}
		storeInterval = fsi
	}
	return nil
}
