package main

import (
	"metrics"
	"net/http"
)

var metricStorage = metrics.NewMemStorage()

func update(res http.ResponseWriter, req *http.Request) {

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", update)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
