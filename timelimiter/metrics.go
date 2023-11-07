package timelimiter

import "sync/atomic"

type Metrics interface {
	SuccessCount() int64
	TimeoutCount() int64
	FailureCount() int64

	successIncrement()
	timeoutIncrement()
	failureIncrement()
}

func newMetric() Metrics {
	return &metrics{}
}

type metrics struct {
	successCounter atomic.Int64
	timeoutCounter atomic.Int64
	failureCounter atomic.Int64
}

func (m *metrics) SuccessCount() int64 {
	return m.successCounter.Load()
}

func (m *metrics) TimeoutCount() int64 {
	return m.timeoutCounter.Load()
}

func (m *metrics) FailureCount() int64 {
	return m.failureCounter.Load()
}

func (m *metrics) successIncrement() {
	m.successCounter.Add(1)
}

func (m *metrics) timeoutIncrement() {
	m.timeoutCounter.Add(1)
}

func (m *metrics) failureIncrement() {
	m.failureCounter.Add(1)
}
