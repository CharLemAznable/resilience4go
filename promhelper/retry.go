package promhelper

import (
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	rtKindSuccessfulWithoutRetry = "successful_without_retry"
	rtKindSuccessfulWithRetry    = "successful_with_retry"
	rtKindFailedWithoutRetry     = "failed_without_retry"
	rtKindFailedWithRetry        = "failed_with_retry"
)

func RetryRegistry(entry retry.Retry) (RegisterFn, UnregisterFn) {
	numberOfSuccessfulCallsWithoutRetryAttemptCounter := prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name:        "resilience4go_retry_calls",
			Help:        "The number of calls with/without retry",
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: rtKindSuccessfulWithoutRetry},
		},
		func() float64 {
			return float64(entry.Metrics().NumberOfSuccessfulCallsWithoutRetryAttempt())
		},
	)
	numberOfSuccessfulCallsWithRetryAttemptCounter := prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name:        "resilience4go_retry_calls",
			Help:        "The number of calls with/without retry",
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: rtKindSuccessfulWithRetry},
		},
		func() float64 {
			return float64(entry.Metrics().NumberOfSuccessfulCallsWithRetryAttempt())
		},
	)
	numberOfFailedCallsWithoutRetryAttemptCounter := prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name:        "resilience4go_retry_calls",
			Help:        "The number of calls with/without retry",
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: rtKindFailedWithoutRetry},
		},
		func() float64 {
			return float64(entry.Metrics().NumberOfFailedCallsWithoutRetryAttempt())
		},
	)
	numberOfFailedCallsWithRetryAttempt := prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name:        "resilience4go_retry_calls",
			Help:        "The number of calls with/without retry",
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: rtKindFailedWithRetry},
		},
		func() float64 {
			return float64(entry.Metrics().NumberOfFailedCallsWithRetryAttempt())
		},
	)
	collectors := []prometheus.Collector{
		numberOfSuccessfulCallsWithoutRetryAttemptCounter,
		numberOfSuccessfulCallsWithRetryAttemptCounter,
		numberOfFailedCallsWithoutRetryAttemptCounter,
		numberOfFailedCallsWithRetryAttempt,
	}
	return buildRegisterFn(collectors...), buildUnregisterFn(collectors...)
}
