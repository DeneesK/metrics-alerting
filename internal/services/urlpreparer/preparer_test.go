package urlpreparer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareURL(t *testing.T) {
	type args struct {
		metricType string
		metricName string
	}
	type want struct {
		result string
	}
	tests := []struct {
		name  string
		args  args
		value float64
		want  want
	}{
		{
			name: "url preparer positive test #1",
			args: args{
				metricName: "metric",
				metricType: "gauge",
			},
			value: 1.5,
			want:  want{"http://localhost:8080/update/gauge/metric/1.5",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := PrepareURL("localhost:8080", test.args.metricType, test.args.metricName, test.value)
			assert.NoError(t, err)
			assert.Equal(t, res, test.want.result)
		})
	}
}
