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
		metrics, err := urlparser.ParseMetricUrl(req.URL.String())
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		memstorage.SaveMetrics(metrics, &metricsStorage)
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
