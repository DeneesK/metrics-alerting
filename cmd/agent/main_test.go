package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sendReport(t *testing.T) {
	type args struct {
		url         string
		contentType string
	}
	tests := []struct {
		name           string
		args           args
		wantResp       string
		wantContenType string
		wantCode       int
	}{
		{
			name: "positive test #1",
			args: args{
				url:         "update/gauge/metric/1.5",
				contentType: "text/plain",
			},
			wantResp:       "/update/gauge/metric/1.5",
			wantContenType: "text/plain; charset=utf-8",
			wantCode:       200,
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", r.Header.Get("Content-Type"))
			w.Write([]byte(r.URL.String()))
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer ts.Close()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url, _ := url.JoinPath(ts.URL, test.args.url)
			resp, err := sendReport(url, test.args.contentType)
			assert.NoError(t, err)
			assert.Equal(t, resp.Header.Get("Content-Type"), test.wantContenType)
			assert.Equal(t, resp.StatusCode, test.wantCode)
			buff := make([]byte, resp.ContentLength)
			resp.Body.Read(buff)
			respBody := bytes.NewBuffer(buff).String()
			assert.Equal(t, respBody, test.wantResp)
		})
	}
}
