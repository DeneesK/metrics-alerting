package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

var RunAddr string
var ReportInterval int
var PollInterval int

func parseFlags() {
	flag.StringVar(&RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&ReportInterval, "r", 10, "interval of sending metrics to the server")
	flag.IntVar(&PollInterval, "p", 2, "interval of polling metrics from the runtime package")
	flag.Parse()
	if envRunAddr, ok := os.LookupEnv("ADDRESS"); ok {
		RunAddr = envRunAddr
	}
	if envreportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		fri, err := strconv.Atoi(envreportInterval)
		if err != nil {
			log.Fatal("The value of the REPORT_INTERVAL environment variable is not a integer.")
		}
		ReportInterval = fri
	}
	if envpolltInterval, ok := os.LookupEnv("POLL_INTERVAL "); ok {
		fpi, err := strconv.Atoi(envpolltInterval)
		if err != nil {
			log.Fatal("The value of the POLL_INTERVAL environment variable is not a integer.")
		}
		PollInterval = fpi
	}
}
