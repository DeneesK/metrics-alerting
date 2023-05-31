package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Conf struct {
	runAddr       string
	logLevel      string
	storeInterval int
	filePath      string
	isRestore     bool
}

func parseFlags() (Conf, error) {
	var conf Conf
	flag.StringVar(&conf.runAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&conf.logLevel, "l", "info", "log level")
	flag.StringVar(&conf.filePath, "f", "tmp/metrics-db.json", "path to store file")
	flag.BoolVar(&conf.isRestore, "r", true, "load saved data")
	flag.IntVar(&conf.storeInterval, "i", 5, "interval of storing data on disk")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		conf.runAddr = envRunAddr
	}
	if envRunAddr := os.Getenv("LOG_LEVEL"); envRunAddr != "" {
		conf.runAddr = envRunAddr
	}
	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		conf.filePath = envFilePath
	}
	if envIsRestore := os.Getenv("RESTORE"); envIsRestore != "" {
		envIsRestore, err := strconv.ParseBool(envIsRestore)
		if err != nil {
			return Conf{}, fmt.Errorf("value of RESTORE is not a boolean: %w", err)
		}
		conf.isRestore = envIsRestore
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		fsi, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			return Conf{}, fmt.Errorf("value of STORE_INTERVAL is not a integer: %w", err)
		}
		conf.storeInterval = fsi
	}
	return conf, nil
}
