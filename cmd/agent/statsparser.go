package main

import (
	"runtime"

	"github.com/fatih/structs"
)

var stats = []string{"Alloc", "BuckHashSys", "Frees", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc"}

func ParseNeededStats(ms runtime.MemStats) map[string]float64 {
	m := structs.Map(ms)
	mapstats := make(map[string]float64)
	for _, stat := range stats {
		mapstats[stat] = float64(m[stat].(uint64))
	}
	mapstats["GCCPUFraction"] = m["GCCPUFraction"].(float64)
	mapstats["NumForcedGC"] = float64(m["NumForcedGC"].(uint32))
	mapstats["NumGC"] = float64(m["NumGC"].(uint32))
	return mapstats
}
