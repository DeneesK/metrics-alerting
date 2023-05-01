package main

import (
	"flag"
)

var flagRunAddr string
var flagreportInterval int
var flagpolltInterval int

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagreportInterval, "r", 10, "override reportInterval - the frequency of sending metrics to the server")
	flag.IntVar(&flagpolltInterval, "p", 2, "override pollInterval - the frequency of polling metrics from the runtime package")
	flag.Parse()
}
