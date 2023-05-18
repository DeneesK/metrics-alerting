package main

import (
	"flag"
	"os"
)

var runAddr string
var logLevel string

func parseFlags() {
	flag.StringVar(&runAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&logLevel, "l", "info", "log level")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		runAddr = envRunAddr
	}
	if envRunAddr := os.Getenv("LOG_LEVEL"); envRunAddr != "" {
		runAddr = envRunAddr
	}
}
