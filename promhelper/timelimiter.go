package promhelper

import (
	"github.com/CharLemAznable/resilience4go/timelimiter"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	tlKinkSuccessful = "successful"
	tlKinkTimeout    = "timeout"
	tlKinkFailed     = "failed"
)

func TimeLimiterRegistry(entry timelimiter.TimeLimiter) (RegisterFn, UnregisterFn) {
	successCounter := prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name:        "resilience4go_timelimiter_calls",
			Help:        "The number of successful calls",
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: tlKinkSuccessful},
		},
		func() float64 {
			return float64(entry.Metrics().SuccessCount())
		},
	)
	timeoutCounter := prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name:        "resilience4go_timelimiter_calls",
			Help:        "The number of timed out calls",
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: tlKinkTimeout},
		},
		func() float64 {
			return float64(entry.Metrics().TimeoutCount())
		},
	)
	failureCounter := prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name:        "resilience4go_timelimiter_calls",
			Help:        "The number of failed calls",
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: tlKinkFailed},
		},
		func() float64 {
			return float64(entry.Metrics().FailureCount())
		},
	)
	collectors := []prometheus.Collector{successCounter, timeoutCounter, failureCounter}
	return buildRegisterFn(collectors...), buildUnregisterFn(collectors...)
}
