package report

import (
	"fmt"
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

type Details struct{}

type Report struct {
	results  chan RequestResult
	finished chan bool
	details  Details
}

func New(results chan RequestResult) Report {
	return Report{
		results:  results,
		finished: make(chan bool, 1),
	}
}

func (r Report) Run() {
	for r := range r.results {
		fmt.Printf("result received: %v\n", r)
	}

	r.finished <- true
}

func (r Report) Finished() <-chan bool {
	return r.finished
}

func (r Report) Details() Details {
	return r.details
}
