package urlpreparer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareURL(t *testing.T) {
	type args struct {
		metricType string
		metricName string
		value      float64
	}
	type want struct {
		result string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "url preparer positive test #1",
			args: args{
				metricName: "metric",
				metricType: "gauge",
				value:      1.5},
			want: want{"http://localhost:8080/update/gauge/metric/1.500000"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := PrepareURL("localhost:8080", test.args.metricType, test.args.metricName, test.args.value)
			fmt.Println(err)
			assert.Equal(t, res, test.want.result)
		})
	}
}