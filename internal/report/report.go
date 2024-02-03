package report

import (
	"math"
	"time"
)

type RequestResult struct {
	DNSDuration        time.Duration
	ConnectionDuration time.Duration
	ResponseDuration   time.Duration
	RequestDuration    time.Duration
	DelayDuration      time.Duration
	TotalDuration      time.Duration
	ResponseCode       int
}

type Details struct {
	ExpectedNumberOfRequests int
	ActualNumberOfRequests   int

	MinDNSLookup     float64
	MaxDNSLookup     float64
	AverageDNSLookup float64

	MinConnection     float64
	MaxConnection     float64
	AverageConnection float64

	MinResponse     float64
	MaxResponse     float64
	AverageResponse float64

	MinRequest     float64
	MaxRequest     float64
	AverageRequest float64

	MinDelay     float64
	MaxDelay     float64
	AverageDelay float64

	MinTotalLatency     float64
	MaxTotalLatency     float64
	AverageTotalLatency float64
}

type Report struct {
	results  chan RequestResult
	finished chan bool
	details  Details
	Test     int
}

func New(expectedRequests int, results chan RequestResult) *Report {
	return &Report{
		results:  results,
		finished: make(chan bool, 1),
		details: Details{
			MinDNSLookup:    math.MaxFloat64,
			MaxDNSLookup:    math.SmallestNonzeroFloat64,
			MinConnection:   math.MaxFloat64,
			MaxConnection:   math.SmallestNonzeroFloat64,
			MinResponse:     math.MaxFloat64,
			MaxResponse:     math.SmallestNonzeroFloat64,
			MinRequest:      math.MaxFloat64,
			MaxRequest:      math.SmallestNonzeroFloat64,
			MinDelay:        math.MaxFloat64,
			MaxDelay:        math.SmallestNonzeroFloat64,
			MinTotalLatency: math.MaxFloat64,
			MaxTotalLatency: math.SmallestNonzeroFloat64,
		},
	}
}

// Run will block until the results channel is closed.
func (r *Report) Run() {
	for res := range r.results {
		r.details.ActualNumberOfRequests++

		r.details.MinDNSLookup = min(r.details.MinDNSLookup, res.DNSDuration.Seconds())
		r.details.MaxDNSLookup = max(r.details.MaxDNSLookup, res.DNSDuration.Seconds())
		r.details.AverageDNSLookup = calculateAverage(
			r.details.AverageDNSLookup,
			res.DNSDuration.Seconds(),
			r.details.ActualNumberOfRequests,
		)

		r.details.MinConnection = min(r.details.MinConnection, res.ConnectionDuration.Seconds())
		r.details.MaxConnection = max(r.details.MaxConnection, res.ConnectionDuration.Seconds())
		r.details.AverageConnection = calculateAverage(
			r.details.AverageConnection,
			res.ConnectionDuration.Seconds(),
			r.details.ActualNumberOfRequests,
		)

		r.details.MinResponse = min(r.details.MinResponse, res.ResponseDuration.Seconds())
		r.details.MaxResponse = max(r.details.MaxResponse, res.ResponseDuration.Seconds())
		r.details.AverageResponse = calculateAverage(
			r.details.AverageResponse,
			res.ResponseDuration.Seconds(),
			r.details.ActualNumberOfRequests,
		)

		r.details.MinRequest = min(r.details.MinRequest, res.RequestDuration.Seconds())
		r.details.MaxRequest = max(r.details.MaxRequest, res.RequestDuration.Seconds())
		r.details.AverageRequest = calculateAverage(
			r.details.AverageRequest,
			res.RequestDuration.Seconds(),
			r.details.ActualNumberOfRequests,
		)

		r.details.MinDelay = min(r.details.MinDelay, res.DelayDuration.Seconds())
		r.details.MaxDelay = max(r.details.MaxDelay, res.DelayDuration.Seconds())
		r.details.AverageDelay = calculateAverage(
			r.details.AverageDelay,
			res.DelayDuration.Seconds(),
			r.details.ActualNumberOfRequests,
		)

		r.details.MinTotalLatency = min(r.details.MinTotalLatency, res.TotalDuration.Seconds())
		r.details.MaxTotalLatency = max(r.details.MaxTotalLatency, res.TotalDuration.Seconds())
		r.details.AverageTotalLatency = calculateAverage(
			r.details.AverageTotalLatency,
			res.TotalDuration.Seconds(),
			r.details.ActualNumberOfRequests,
		)
	}

	r.finished <- true
}

// Finished returns a channel that will be closed when the reporter is finished.
func (r *Report) Finished() <-chan bool {
	return r.finished
}

// Details returns the details of the report.
func (r *Report) Details() Details {
	return r.details
}

func calculateAverage(oldAverage, newValue float64, count int) float64 {
	return (oldAverage*(float64(count)-1) + newValue) / float64(count)
}
