package main

import (
	"fmt"
	"net/http"

	"github.com/DeneesK/metrics-alerting/cmd/server/memstorage"
	"github.com/DeneesK/metrics-alerting/cmd/server/urlparser"
)

var metricsStorage = memstorage.NewMemStorage()

func update(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		err := memstorage.SaveMetrics(req.URL.String(), &metricsStorage)
		if err != nil {
			switch err {
			case urlparser.ErrConverValue:
				res.WriteHeader(http.StatusBadRequest)
			case urlparser.ErrMetricType:
				res.WriteHeader(http.StatusBadRequest)
			case urlparser.ErrEmptyMetricName:
				res.WriteHeader(http.StatusNotFound)
			}
			return
		}
		res.WriteHeader(http.StatusOK)
		fmt.Printf("gauge: %v\ncounter: %v\n", metricsStorage.Gauge, metricsStorage.Counter)
		return
	}
	res.WriteHeader(http.StatusMethodNotAllowed)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", update)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
