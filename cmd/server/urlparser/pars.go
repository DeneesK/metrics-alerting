package urlparser

import (
	"errors"
	"strconv"
	"strings"
)

var ErrConverValue = errors.New("invalid type of value")
var ErrMetricType = errors.New("invalid type of metric")
var ErrEmptyMetricName = errors.New("missed metric name")

type MetricPayload struct {
	MetricType   string
	MetricName   string
	CaugeValue   float32
	CounterValue int
}

const (
	counterMetric string = "counter"
	gaugeMetric   string = "gauge"
)

func ParseMetricURL(u string) (MetricPayload, error) {
	mp := MetricPayload{}
	arr := strings.Split(u, "/")[2:]
	switch arr[0] {
	case counterMetric:
		mp.MetricType = counterMetric
		if len(arr) < 3 {
			return mp, ErrEmptyMetricName
		}
		n, err := strconv.Atoi(arr[2])
		if err != nil {
			return mp, ErrConverValue
		}
		mp.CounterValue = n
		if len(arr) < 3 {
			return mp, ErrEmptyMetricName
		}
		mp.MetricName = arr[1]
	case gaugeMetric:
		if len(arr) < 3 {
			return mp, ErrEmptyMetricName
		}
		mp.MetricType = gaugeMetric
		f, err := strconv.ParseFloat(arr[2], 32)
		if err != nil {
			return mp, ErrConverValue
		}
		mp.CaugeValue = float32(f)
		mp.MetricName = arr[1]
	default:
		return mp, ErrMetricType
	}
	return mp, nil
}
