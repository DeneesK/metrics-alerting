package main

import (
	"flag"
	"os"
)

var RunAddr string

func parseFlags() {
	flag.StringVar(&RunAddr, "a", ":8080", "address and port to run server")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		RunAddr = envRunAddr
	}
}
