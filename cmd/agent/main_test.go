package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/DeneesK/metrics-alerting/internal/models"
	"github.com/levigross/grequests"
	"github.com/stretchr/testify/assert"
)

func Test_postReport(t *testing.T) {
	v := 10.5
	metrics := make([]models.Metrics, 0)
	type args struct {
		metrics []models.Metrics
	}
	tests := []struct {
		name            string
		args            args
		wantContentType string
		wantCode        int
	}{
		{
			name: "positive test #1",
			args: args{
				metrics: append(metrics, models.Metrics{ID: "PollCount", MType: "gauge", Value: &v}),
			},
			wantContentType: "application/json",
			wantCode:        200,
		},
	}
	ro := grequests.RequestOptions{Headers: map[string]string{"Content-Type": "application/json"}}
	session := grequests.NewSession(&ro)
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost {
					assert.Equal(t, r.Header.Get("Content-Type"), test.wantContentType)
					w.WriteHeader(http.StatusOK)
				}
				w.WriteHeader(http.StatusMethodNotAllowed)
			}))
			url, err := url.JoinPath(ts.URL, "update")
			assert.NoError(t, err)
			statusCode, err := sendBanch(session, url, test.args.metrics)
			assert.NoError(t, err)
			assert.Equal(t, statusCode, test.wantCode)
			ts.Close()
		})
	}
}
