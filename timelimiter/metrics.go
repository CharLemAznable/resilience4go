package timelimiter

import "sync/atomic"

type Metrics interface {
	SuccessCount() int64
	TimeoutCount() int64
	PanicCount() int64

	successIncrement()
	timeoutIncrement()
	panicIncrement()
}

func newMetrics() Metrics {
	return &metrics{}
}

type metrics struct {
	successCounter atomic.Int64
	timeoutCounter atomic.Int64
	panicCounter   atomic.Int64
}

func (m *metrics) SuccessCount() int64 {
	return m.successCounter.Load()
}

func (m *metrics) TimeoutCount() int64 {
	return m.timeoutCounter.Load()
}

func (m *metrics) PanicCount() int64 {
	return m.panicCounter.Load()
}

func (m *metrics) successIncrement() {
	m.successCounter.Add(1)
}

func (m *metrics) timeoutIncrement() {
	m.timeoutCounter.Add(1)
}

func (m *metrics) panicIncrement() {
	m.panicCounter.Add(1)
}
