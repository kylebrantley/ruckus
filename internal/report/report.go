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
}

type Report struct {
	results  chan RequestResult
	finished chan bool
	details  Details
	Test     int
}

func New(results chan RequestResult) *Report {
	return &Report{
		results:  results,
		finished: make(chan bool, 1),
	}
}

func (r *Report) Run() {
	for res := range r.results {
		_ = res
		r.details.ActualNumberOfRequests++
	}

	r.finished <- true
}

func (r *Report) Finished() <-chan bool {
	return r.finished
}

func (r *Report) Details() Details {
	return r.details
}
