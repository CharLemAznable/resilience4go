package promhelper

import (
	"github.com/CharLemAznable/resilience4go/timelimiter"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	tlKinkSuccessful = "successful"
	tlKinkTimeout    = "timeout"
	tlKinkPanicked   = "panicked"

	timelimiterCallsName = "resilience4go_timelimiter_calls"
	timelimiterCallsHelp = "The number of successful/timeout/panicked calls"
)

func TimeLimiterRegistry(entry timelimiter.TimeLimiter) (RegisterFn, UnregisterFn) {
	successCounter := prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name:        timelimiterCallsName,
			Help:        timelimiterCallsHelp,
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: tlKinkSuccessful},
		},
		func() float64 {
			return float64(entry.Metrics().SuccessCount())
		},
	)
	timeoutCounter := prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name:        timelimiterCallsName,
			Help:        timelimiterCallsHelp,
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: tlKinkTimeout},
		},
		func() float64 {
			return float64(entry.Metrics().TimeoutCount())
		},
	)
	panicCounter := prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name:        timelimiterCallsName,
			Help:        timelimiterCallsHelp,
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: tlKinkPanicked},
		},
		func() float64 {
			return float64(entry.Metrics().PanicCount())
		},
	)
	collectors := []prometheus.Collector{successCounter, timeoutCounter, panicCounter}
	return buildRegisterFn(collectors...), buildUnregisterFn(collectors...)
}
