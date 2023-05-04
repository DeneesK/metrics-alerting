package urlpreparer

import (
	"fmt"
	"net/url"
)

func PrepareURL(basePath string, metricType string, metricName string, value float64) (string, error) {
	v := fmt.Sprintf("%f", value)
	u, err := url.JoinPath("http://", basePath, "update", metricType, metricName, v)
	if err != nil {
		return "", fmt.Errorf("PrepareURL failed: %v", err)
	}
	return u, nil
}
