package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/levigross/grequests"
	"github.com/stretchr/testify/assert"
)

func Test_postReport(t *testing.T) {
	type args struct {
		contentType string
		metricType  string
		metricName  string
		delta       uint64
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
				contentType: "application/json",
				metricType:  "counter",
				metricName:  "metric",
				delta:       10,
			},
			wantContentType: "application/json",
			wantCode:        200,
		},
	}
	ro := grequests.RequestOptions{Headers: map[string]string{"Content-Type": "application/json"}}
	session := grequests.NewSession(&ro)
	for _, test := range tests {
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
			statusCode, err := send(session, url, test.args.metricType, test.args.metricName, test.args.delta)
			assert.NoError(t, err)
			assert.Equal(t, statusCode, test.wantCode)
			ts.Close()
		})
	}
}
