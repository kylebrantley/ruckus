package report

import (
	"math"
)

type Metric struct {
	Min     float64
	Max     float64
	Average float64
}

func NewMetric() Metric {
	return Metric{
		Min:     math.MaxFloat64,
		Max:     math.SmallestNonzeroFloat64,
		Average: 0,
	}
}

func (s *Metric) Update(value float64, count int) {
	s.Min = math.Min(s.Min, value)
	s.Max = math.Max(s.Max, value)
	s.Average += (value - s.Average) / float64(count)
}
