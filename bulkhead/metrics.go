package bulkhead

import "sync/atomic"

type Metrics interface {
	MaxAllowedConcurrentCalls() int64
	AvailableConcurrentCalls() int64

	acquire(n int64)
	release(n int64)
}

func newMetrics(maxConcurrentCalls int64) Metrics {
	m := &metrics{
		maxConcurrentCalls: maxConcurrentCalls,
	}
	m.availablePermits.Store(maxConcurrentCalls)
	return m
}

type metrics struct {
	maxConcurrentCalls int64
	availablePermits   atomic.Int64
}

func (m *metrics) MaxAllowedConcurrentCalls() int64 {
	return m.maxConcurrentCalls
}

func (m *metrics) AvailableConcurrentCalls() int64 {
	return m.availablePermits.Load()
}

func (m *metrics) acquire(n int64) {
	m.availablePermits.Add(-n)
}

func (m *metrics) release(n int64) {
	m.availablePermits.Add(n)
}
