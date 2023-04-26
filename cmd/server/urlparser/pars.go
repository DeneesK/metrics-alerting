package urlparser

import (
	"errors"
	"strconv"
	"strings"
)

type MetricPayload struct {
	MetricType   string
	MetricName   string
	CaugeValue   float32
	CounterValue int
}

func ParseMetricUrl(u string) (MetricPayload, error) {
	mp := MetricPayload{}
	arr := strings.Split(u, "/")[2:]
	switch arr[0] {
	case "counter":
		mp.MetricType = "counter"
		n, err := strconv.Atoi(arr[2])
		if err != nil {
			return mp, err
		}
		mp.CounterValue = n
		mp.MetricName = arr[1]
	case "gauge":
		mp.MetricType = "gauge"
		f, err := strconv.ParseFloat(arr[2], 32)
		if err != nil {
			return mp, err
		}
		mp.CaugeValue = float32(f)
		mp.MetricName = arr[1]
	default:
		return mp, errors.New("invalid type of metric")
	}
	return mp, nil
}
