package promhelper

import (
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/prometheus/client_golang/prometheus"
)

func RetryCollectors(entry retry.Retry) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewCounterFunc(
			prometheus.CounterOpts{
				Name:        "resilience4go_retry_calls",
				Help:        "The number of successful calls without a retry attempt",
				ConstLabels: prometheus.Labels{"name": entry.Name(), "kind": "successful_without_retry"},
			},
			func() float64 {
				return float64(entry.Metrics().NumberOfSuccessfulCallsWithoutRetryAttempt())
			},
		),
		prometheus.NewCounterFunc(
			prometheus.CounterOpts{
				Name:        "resilience4go_retry_calls",
				Help:        "The number of successful calls after a retry attempt",
				ConstLabels: prometheus.Labels{"name": entry.Name(), "kind": "successful_with_retry"},
			},
			func() float64 {
				return float64(entry.Metrics().NumberOfSuccessfulCallsWithRetryAttempt())
			},
		),
		prometheus.NewCounterFunc(
			prometheus.CounterOpts{
				Name:        "resilience4go_retry_calls",
				Help:        "The number of failed calls without a retry attempt",
				ConstLabels: prometheus.Labels{"name": entry.Name(), "kind": "failed_without_retry"},
			},
			func() float64 {
				return float64(entry.Metrics().NumberOfFailedCallsWithoutRetryAttempt())
			},
		),
		prometheus.NewCounterFunc(
			prometheus.CounterOpts{
				Name:        "resilience4go_retry_calls",
				Help:        "The number of failed calls after a retry attempt",
				ConstLabels: prometheus.Labels{"name": entry.Name(), "kind": "failed_with_retry"},
			},
			func() float64 {
				return float64(entry.Metrics().NumberOfFailedCallsWithRetryAttempt())
			},
		),
	}
}
