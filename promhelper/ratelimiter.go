package promhelper

import (
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/prometheus/client_golang/prometheus"
)

func RateLimiterRegistry(entry ratelimiter.RateLimiter) (RegisterFn, UnregisterFn) {
	numberOfWaitingThreadsGauge := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "resilience4go_ratelimiter_waiting_threads",
			Help:        "The number of waiting threads",
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name()},
		},
		func() float64 {
			return float64(entry.Metrics().NumberOfWaitingThreads())
		},
	)
	availablePermissionsGauge := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "resilience4go_ratelimiter_available_permissions",
			Help:        "The number of available permissions",
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name()},
		},
		func() float64 {
			return float64(entry.Metrics().AvailablePermissions())
		},
	)
	collectors := []prometheus.Collector{
		numberOfWaitingThreadsGauge,
		availablePermissionsGauge,
	}
	return buildRegisterFn(collectors...), buildUnregisterFn(collectors...)
}
