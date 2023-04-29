package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DeneesK/metrics-alerting/cmd/server/memstorage"
	"github.com/stretchr/testify/assert"
)

func Test_update(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "positive test #1",
			args: "/update/counter/metric/1",
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test: empty value #1",
			args: "/update/counter/metric/",
			want: want{
				code:        400,
				contentType: "",
			},
		},
		{
			name: "negative test: missing metric name #1",
			args: "/update/counter/",
			want: want{
				code:        404,
				contentType: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ms := memstorage.NewMemStorage()
			request := httptest.NewRequest(http.MethodPost, test.args, nil)
			w := httptest.NewRecorder()
			update(&ms)(w, request)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, res.StatusCode, test.want.code)
			assert.Equal(t, res.Header.Get("Content-Type"), test.want.contentType)
		})
	}
}
