package timelimiter

import "sync/atomic"

type Metrics interface {
	SuccessCount() uint64
	TimeoutCount() uint64
	PanicCount() uint64
}

func newMetrics() *metrics {
	return &metrics{}
}

type metrics struct {
	successCounter atomic.Uint64
	timeoutCounter atomic.Uint64
	panicCounter   atomic.Uint64
}

func (m *metrics) SuccessCount() uint64 {
	return m.successCounter.Load()
}

func (m *metrics) TimeoutCount() uint64 {
	return m.timeoutCounter.Load()
}

func (m *metrics) PanicCount() uint64 {
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
