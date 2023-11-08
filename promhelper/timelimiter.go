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

func TimeLimiterCollectors(entry timelimiter.TimeLimiter) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewCounterFunc(
			prometheus.CounterOpts{
				Name:        "resilience4go_timelimiter_calls",
				Help:        "The number of successful calls",
				ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: tlKinkSuccessful},
			},
			func() float64 {
				return float64(entry.Metrics().SuccessCount())
			},
		),
		prometheus.NewCounterFunc(
			prometheus.CounterOpts{
				Name:        "resilience4go_timelimiter_calls",
				Help:        "The number of timed out calls",
				ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: tlKinkTimeout},
			},
			func() float64 {
				return float64(entry.Metrics().TimeoutCount())
			},
		),
		prometheus.NewCounterFunc(
			prometheus.CounterOpts{
				Name:        "resilience4go_timelimiter_calls",
				Help:        "The number of failed calls",
				ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: tlKinkFailed},
			},
			func() float64 {
				return float64(entry.Metrics().FailureCount())
			},
		),
	}
}
