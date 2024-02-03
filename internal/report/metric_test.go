package report

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricUpdate(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   Metric
	}{
		{
			name:   "single value",
			values: []float64{5.0},
			want:   Metric{Min: 5.0, Max: 5.0, Average: 5.0},
		},
		{
			name:   "multiple values",
			values: []float64{5.0, 10.0, 15.0},
			want:   Metric{Min: 5.0, Max: 15.0, Average: 10.0},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				m := NewMetric()
				for i, value := range tt.values {
					m.Update(value, i+1)
				}

				assert.Equal(t, tt.want, m)
			},
		)
	}
}
