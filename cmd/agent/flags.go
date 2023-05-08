package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var runAddr string
var reportingInterval int
var pollingInterval int

func parseFlags() error {
	flag.StringVar(&runAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&reportingInterval, "r", 10, "interval of sending metrics to the server")
	flag.IntVar(&pollingInterval, "p", 2, "interval of polling metrics from the runtime package")
	flag.Parse()
	if envRunAddr, ok := os.LookupEnv("ADDRESS"); ok {
		runAddr = envRunAddr
	}
	if envreportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		fri, err := strconv.Atoi(envreportInterval)
		if err != nil {
			err = fmt.Errorf("value of REPORT_INTERVAL is not a integer: %w", err)
			return err
		}
		reportingInterval = fri
	}
	if envpollInterval, ok := os.LookupEnv("POLL_INTERVAL "); ok {
		fpi, err := strconv.Atoi(envpollInterval)
		if err != nil {
			err = fmt.Errorf("value of POLL_INTERVAL is not a integer: %w", err)
			return err
		}
		pollingInterval = fpi
	}
	return nil
}
