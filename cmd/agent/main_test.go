package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/levigross/grequests"
	"github.com/stretchr/testify/assert"
)

func Test_sendReport(t *testing.T) {
	type args struct {
		url         string
		contentType string
	}
	tests := []struct {
		name            string
		args            args
		wantURL         string
		wantContentType string
		wantCode        int
	}{
		{
			name: "positive test #1",
			args: args{
				url:         "update/gauge/metric/1.5",
				contentType: "text/plain",
			},
			wantURL:         "/update/gauge/metric/1.5",
			wantContentType: "text/plain",
			wantCode:        200,
		},
	}
	ro := grequests.RequestOptions{Headers: map[string]string{"Content-Type": "text/plain"}}
	session := grequests.NewSession(&ro)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost {
					assert.Equal(t, r.URL.String(), test.wantURL)
					assert.Equal(t, r.Header.Get("Content-Type"), test.wantContentType)
					w.WriteHeader(http.StatusOK)
				}
				w.WriteHeader(http.StatusMethodNotAllowed)
			}))
			url, _ := url.JoinPath(ts.URL, test.args.url)
			resp, err := sendReport(session, url)
			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode, test.wantCode)
			ts.Close()
		})
	}
}
