package promhelper

import (
	"github.com/CharLemAznable/resilience4go/timelimiter"
	"github.com/prometheus/client_golang/prometheus"
)

func TimeLimiterCollectors(entry timelimiter.TimeLimiter) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewCounterFunc(
			prometheus.CounterOpts{
				Name:        "resilience4go_timelimiter_calls",
				Help:        "The number of successful calls",
				ConstLabels: prometheus.Labels{"name": entry.Name(), "kind": "successful"},
			},
			func() float64 {
				return float64(entry.Metrics().SuccessCount())
			},
		),
		prometheus.NewCounterFunc(
			prometheus.CounterOpts{
				Name:        "resilience4go_timelimiter_calls",
				Help:        "The number of timed out calls",
				ConstLabels: prometheus.Labels{"name": entry.Name(), "kind": "timeout"},
			},
			func() float64 {
				return float64(entry.Metrics().TimeoutCount())
			},
		),
		prometheus.NewCounterFunc(
			prometheus.CounterOpts{
				Name:        "resilience4go_timelimiter_calls",
				Help:        "The number of failed calls",
				ConstLabels: prometheus.Labels{"name": entry.Name(), "kind": "failed"},
			},
			func() float64 {
				return float64(entry.Metrics().FailureCount())
			},
		),
	}
}
