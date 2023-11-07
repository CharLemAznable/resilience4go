package retry

import "sync/atomic"

type Metrics interface {
	NumberOfSuccessfulCallsWithoutRetryAttempt() int64
	NumberOfSuccessfulCallsWithRetryAttempt() int64
	NumberOfFailedCallsWithoutRetryAttempt() int64
	NumberOfFailedCallsWithRetryAttempt() int64

	successfulCallsWithoutRetryAttemptIncrement()
	successfulCallsWithRetryAttemptIncrement()
	failedCallsWithoutRetryAttemptIncrement()
	failedCallsWithRetryAttemptIncrement()
}

func newMetrics() Metrics {
	return &metrics{}
}

type metrics struct {
	succeededWithoutRetryCounter atomic.Int64
	succeededAfterRetryCounter   atomic.Int64
	failedWithoutRetryCounter    atomic.Int64
	failedAfterRetryCounter      atomic.Int64
}

func (m *metrics) NumberOfSuccessfulCallsWithoutRetryAttempt() int64 {
	return m.succeededWithoutRetryCounter.Load()
}

func (m *metrics) NumberOfSuccessfulCallsWithRetryAttempt() int64 {
	return m.succeededAfterRetryCounter.Load()
}

func (m *metrics) NumberOfFailedCallsWithoutRetryAttempt() int64 {
	return m.failedWithoutRetryCounter.Load()
}

func (m *metrics) NumberOfFailedCallsWithRetryAttempt() int64 {
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
