package report

import (
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

	DNSLookup    Metric
	Connection   Metric
	Response     Metric
	Request      Metric
	Delay        Metric
	TotalLatency Metric
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
			ExpectedNumberOfRequests: expectedRequests,
			DNSLookup:                NewMetric(),
			Connection:               NewMetric(),
			Response:                 NewMetric(),
			Request:                  NewMetric(),
			Delay:                    NewMetric(),
			TotalLatency:             NewMetric(),
		},
	}
}

// Run will block until the results channel is closed.
func (r *Report) Run() {
	for res := range r.results {
		r.details.ActualNumberOfRequests++

		r.details.DNSLookup.Update(res.DNSDuration.Seconds(), r.details.ActualNumberOfRequests)
		r.details.Connection.Update(res.ConnectionDuration.Seconds(), r.details.ActualNumberOfRequests)
		r.details.Response.Update(res.ResponseDuration.Seconds(), r.details.ActualNumberOfRequests)
		r.details.Request.Update(res.RequestDuration.Seconds(), r.details.ActualNumberOfRequests)
		r.details.Delay.Update(res.DelayDuration.Seconds(), r.details.ActualNumberOfRequests)
		r.details.TotalLatency.Update(res.TotalDuration.Seconds(), r.details.ActualNumberOfRequests)
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
