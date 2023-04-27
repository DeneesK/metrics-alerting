package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DeneesK/metrics-alerting/cmd/server/memstorage"
	"github.com/DeneesK/metrics-alerting/cmd/server/urlparser"
)

func update(storage *memstorage.MemStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			err := memstorage.SaveMetrics(req.URL.String(), storage)
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
			fmt.Printf("gauge: %v\ncounter: %v\n", storage.Gauge, storage.Counter)
			return
		}
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	mux := http.NewServeMux()
	metricsStorage := memstorage.NewMemStorage()
	mux.HandleFunc("/update/", update(&metricsStorage))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
