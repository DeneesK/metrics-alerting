package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DeneesK/metrics-alerting/cmd/server/memstorage"
	"github.com/DeneesK/metrics-alerting/cmd/server/urlparser"
)

func update(storage memstorage.Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			err := storage.StoreMetrics(req.URL.String())
			if err != nil {
				switch err {
				case urlparser.ErrConvertValue:
					res.WriteHeader(http.StatusBadRequest)
				case urlparser.ErrMetricType:
					res.WriteHeader(http.StatusBadRequest)
				case urlparser.ErrEmptyMetricName:
					res.WriteHeader(http.StatusNotFound)
				}
				return
			}
			res.Header().Add("Content-Type", "text/plain; charset=utf-8")
			res.WriteHeader(http.StatusOK)
			fmt.Println(storage)
			return
		}
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	mux := http.NewServeMux()
	metricsStorage := memstorage.NewMemStorage()
	mux.HandleFunc("/update/", update(&metricsStorage))
	log.Println("server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
