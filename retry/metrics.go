package retry

import "sync/atomic"

type Metrics interface {
	NumberOfSuccessfulCallsWithoutRetryAttempt() uint64
	NumberOfSuccessfulCallsWithRetryAttempt() uint64
	NumberOfFailedCallsWithoutRetryAttempt() uint64
	NumberOfFailedCallsWithRetryAttempt() uint64

	successfulCallsWithoutRetryAttemptIncrement()
	successfulCallsWithRetryAttemptIncrement()
	failedCallsWithoutRetryAttemptIncrement()
	failedCallsWithRetryAttemptIncrement()
}

func newMetrics() Metrics {
	return &metrics{}
}

type metrics struct {
	succeededWithoutRetryCounter atomic.Uint64
	succeededAfterRetryCounter   atomic.Uint64
	failedWithoutRetryCounter    atomic.Uint64
	failedAfterRetryCounter      atomic.Uint64
}

func (m *metrics) NumberOfSuccessfulCallsWithoutRetryAttempt() uint64 {
	return m.succeededWithoutRetryCounter.Load()
}

func (m *metrics) NumberOfSuccessfulCallsWithRetryAttempt() uint64 {
	return m.succeededAfterRetryCounter.Load()
}

func (m *metrics) NumberOfFailedCallsWithoutRetryAttempt() uint64 {
	return m.failedWithoutRetryCounter.Load()
}

func (m *metrics) NumberOfFailedCallsWithRetryAttempt() uint64 {
	return m.failedAfterRetryCounter.Load()
}

func (m *metrics) successfulCallsWithoutRetryAttemptIncrement() {
	m.succeededWithoutRetryCounter.Add(1)
}

func (m *metrics) successfulCallsWithRetryAttemptIncrement() {
	m.succeededAfterRetryCounter.Add(1)
}

func (m *metrics) failedCallsWithoutRetryAttemptIncrement() {
	m.failedWithoutRetryCounter.Add(1)
}

func (m *metrics) failedCallsWithRetryAttemptIncrement() {
	m.failedAfterRetryCounter.Add(1)
}
