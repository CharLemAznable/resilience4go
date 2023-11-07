package promhelper

import (
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/prometheus/client_golang/prometheus"
)

func RateLimiterCollectors(entry ratelimiter.RateLimiter) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "resilience4go_ratelimiter_waiting_threads",
				Help:        "The number of waiting threads",
				ConstLabels: prometheus.Labels{"name": entry.Name()},
			},
			func() float64 {
				return float64(entry.Metrics().NumberOfWaitingThreads())
			},
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "resilience4go_ratelimiter_available_permissions",
				Help:        "The number of available permissions",
				ConstLabels: prometheus.Labels{"name": entry.Name()},
			},
			func() float64 {
				return float64(entry.Metrics().AvailablePermissions())
			},
		),
	}
}
