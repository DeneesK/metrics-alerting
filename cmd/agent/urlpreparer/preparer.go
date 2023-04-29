package urlpreparer

import (
	"fmt"
	"log"
	"net/url"
	"runtime"

	"github.com/fatih/structs"
)

const (
	basePath string = "http://localhost:8080/update"
)

var stats = []string{"Alloc", "BuckHashSys", "Frees", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc"}

func ParseNeededStats(ms runtime.MemStats) map[string]float32 {
	m := structs.Map(ms)
	mapstats := make(map[string]float32)
	for _, stat := range stats {
		mapstats[stat] = float32(m[stat].(uint64))
	}
	mapstats["GCCPUFraction"] = float32(m["GCCPUFraction"].(float64))
	mapstats["NumForcedGC"] = float32(m["NumForcedGC"].(uint32))
	mapstats["NumGC"] = float32(m["NumGC"].(uint32))
	return mapstats
}

func PrepareURL(metricType string, metricName string, value float32) string {
	v := fmt.Sprintf("%f", value)
	u, err := url.JoinPath(basePath, metricType, metricName, v)
	if err != nil {
		log.Println(err)
		return ""
	}
	return u
}
