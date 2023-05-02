package main

import (
	"flag"
)

var flagRunAddr string
var flagreportInterval int
var flagpolltInterval int

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "http://localhost:8080", "address and port to run server")
	flag.IntVar(&flagreportInterval, "r", 1, "override reportInterval - the frequency of sending metrics to the server")
	flag.IntVar(&flagpolltInterval, "p", 1, "override pollInterval - the frequency of polling metrics from the runtime package")
	flag.Parse()
}
