package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"sync"
	"time"
)

// Arbitrary max channel size to reduce how much memory will be allocated
// TODO: Do some tests and find an optimal value
const maxChannelSize = 100

type result struct {
	DNSDuration        time.Duration
	ConnectionDuration time.Duration
	ResponseDuration   time.Duration
	RequestDuration    time.Duration
	DelayDuration      time.Duration
	TotalDuration      time.Duration
	ResponseCode       int
}

type Executor struct {
	Request          *http.Request
	RequestBody      []byte
	NumberOfRequests int // Total number of requests to make
	NumberOfThreads  int // Maximum number concurrent workers
	RequestTimeout   int

	results chan *result
	stop    chan struct{}
}

func (e *Executor) Init() {
	e.results = make(chan *result, e.NumberOfRequests)
	e.stop = make(chan struct{}, min(e.NumberOfThreads, maxChannelSize)) // nolint:typecheck
}

func (e *Executor) Start() {
	e.Init()
	go func() {
		for r := range e.results {
			fmt.Printf("\nresult: %v\n", r)
		}
	}()
	e.runWorkers()
}

func (e *Executor) Stop() {
	for i := 0; i < e.NumberOfThreads; i++ {
		e.stop <- struct{}{}
	}
	fmt.Printf("shutdown received")
	close(e.results)
	close(e.stop)
}

func (e *Executor) runWorkers() {
	var wg sync.WaitGroup
	wg.Add(e.NumberOfThreads)

	client := &http.Client{
		Timeout: time.Second * time.Duration(e.RequestTimeout),
	}

	for i := 0; i < e.NumberOfThreads; i++ {
		go func() {
			e.runWorker(client, e.NumberOfRequests/e.NumberOfThreads)
			wg.Done()
		}()
	}

	wg.Wait()
}

func (e *Executor) runWorker(client *http.Client, numberOfRequests int) {
	for i := 0; i < numberOfRequests; i++ {
		select {
		case <-e.stop:
			return
		default:
			// TODO: handle error
			_ = e.executeRequest(client)
		}
	}
}

func (e *Executor) executeRequest(client *http.Client) error {
	start := time.Now()

	var dnsStart, connectionStart, responseStart, requestStart, delayStart time.Time
	var dnsDuration, connectionDuration, responseDuration, requestDuration, delayDuration time.Duration

	request, err := cloneRequest(e.Request)
	if err != nil {
		return err
	}

	trace := &httptrace.ClientTrace{
		DNSStart: func(i httptrace.DNSStartInfo) {
			dnsStart = time.Now()
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			dnsDuration = time.Since(dnsStart)
		},
		GetConn: func(h string) {
			connectionStart = time.Now()
		},
		GotConn: func(i httptrace.GotConnInfo) {
			if !i.Reused {
				connectionDuration = time.Since(connectionStart)
			}
			requestStart = time.Now()
		},
		WroteRequest: func(i httptrace.WroteRequestInfo) {
			requestDuration = time.Since(requestStart)
			delayStart = time.Now()
		},
		GotFirstResponseByte: func() {
			delayDuration = time.Since(delayStart)
			responseStart = time.Now()
		},
	}

	request = request.WithContext(httptrace.WithClientTrace(request.Context(), trace))
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("request failed: %v", err)
		// TODO: something
	}

	responseDuration = time.Since(responseStart)
	totalDuration := time.Since(start)
	e.results <- &result{
		DNSDuration:        dnsDuration,
		ConnectionDuration: connectionDuration,
		ResponseDuration:   responseDuration,
		RequestDuration:    requestDuration,
		DelayDuration:      delayDuration,
		TotalDuration:      totalDuration,
		ResponseCode:       response.StatusCode,
	}

	return nil
}

func cloneRequest(r *http.Request) (*http.Request, error) {
	r2 := r.Clone(context.TODO())
	if r.Body == nil {
		return r2, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	if len(body) > 0 {
		r2.Body = io.NopCloser(bytes.NewReader(body))
	}

	return r2, nil
}
