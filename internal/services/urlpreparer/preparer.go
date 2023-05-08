package urlpreparer

import (
	"fmt"
	"net/url"
	"strconv"
)

func PrepareURL(basePath string, metricType string, metricName string, value interface{}) (string, error) {
	switch valueType := value.(type) {
	case uint64:
		v := strconv.FormatUint(value.(uint64), 10)
		return format(basePath, metricType, metricName, v)
	case int64:
		v := strconv.FormatInt(value.(int64), 10)
		return format(basePath, metricType, metricName, v)
	case float64:
		v := strconv.FormatFloat(value.(float64), byte(102), -3, 64)
		return format(basePath, metricType, metricName, v)
	default:
		return "", fmt.Errorf("unable to create url, value must be float64, uint64 or int, given: %v", valueType)
	}
}

func format(basePath, metricType, metricName, v string) (string, error) {
	u, err := url.JoinPath("http://", basePath, "update", metricType, metricName, v)
	if err != nil {
		return "", fmt.Errorf("unable to create url: %v", err)
	}
	return u, nil
}
