package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DeneesK/metrics-alerting/internal/api"
	"github.com/DeneesK/metrics-alerting/internal/models"
	"github.com/DeneesK/metrics-alerting/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var metric_counter int64 = 1

func Test_update(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		args models.Metrics
		want want
	}{
		{
			name: "positive test #1",
			args: models.Metrics{ID: "metric", MType: "counter", Delta: &metric_counter},
			want: want{
				code:        200,
				contentType: "application/json",
			},
		},
	}
	ms := storage.NewMemStorage()
	ts := httptest.NewServer(api.RouterWithoutLogger(&ms))
	defer ts.Close()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := json.Marshal(&test.args)
			require.NoError(t, err)
			buf := bytes.NewBuffer(res)
			request, err := http.NewRequest(http.MethodPost, ts.URL+"/update", buf)
			require.NoError(t, err)
			resp, err := ts.Client().Do(request)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, resp.StatusCode, test.want.code)
			assert.Equal(t, resp.Header.Get("Content-Type"), test.want.contentType)
		})
	}
}
