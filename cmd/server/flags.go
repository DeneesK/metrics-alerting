package main

import (
	"flag"
	"os"
)

var runAddr string

func parseFlags() {
	flag.StringVar(&runAddr, "a", ":8080", "address and port to run server")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		runAddr = envRunAddr
	}
}
