package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Conf struct {
	runAddr           string
	hashKey           string
	reportingInterval int
	pollingInterval   int
}

func parseFlags() (Conf, error) {
	var conf Conf
	flag.StringVar(&conf.runAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&conf.hashKey, "k", "", "the key to calculate hash")
	flag.IntVar(&conf.reportingInterval, "r", 10, "interval of sending metrics to the server")
	flag.IntVar(&conf.pollingInterval, "p", 2, "interval of polling metrics from the runtime package")
	flag.Parse()
	if envRunAddr, ok := os.LookupEnv("ADDRESS"); ok {
		conf.runAddr = envRunAddr
	}
	if envHashKey, ok := os.LookupEnv("KEY"); ok {
		conf.hashKey = envHashKey
	}
	if envreportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		fri, err := strconv.Atoi(envreportInterval)
		if err != nil {
			return Conf{}, fmt.Errorf("value of REPORT_INTERVAL is not a integer: %w", err)
		}
		conf.reportingInterval = fri
	}
	if envpollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		fpi, err := strconv.Atoi(envpollInterval)
		if err != nil {
			return Conf{}, fmt.Errorf("value of POLL_INTERVAL is not a integer: %w", err)
		}
		conf.pollingInterval = fpi
	}
	return conf, nil
}
