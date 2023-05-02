package urlpreparer

import (
	"fmt"
	"log"
	"net/url"
)

func PrepareURL(basePath string, metricType string, metricName string, value float64) string {
	v := fmt.Sprintf("%f", value)
	u, err := url.JoinPath("http://", basePath, "update", metricType, metricName, v)
	if err != nil {
		log.Println(err)
		return ""
	}
	return u
}
