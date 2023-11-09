package promhelper

import (
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/prometheus/client_golang/prometheus"
)

func BulkheadRegistry(entry bulkhead.Bulkhead) (RegisterFn, UnregisterFn) {
	maxAllowedConcurrentCallsGauge := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "resilience4go_bulkhead_max_allowed_concurrent_calls",
			Help:        "The maximum number of available permissions",
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name()},
		},
		func() float64 {
			return float64(entry.Metrics().MaxAllowedConcurrentCalls())
		},
	)
	availableConcurrentCallsGauge := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "resilience4go_bulkhead_available_concurrent_calls",
			Help:        "The number of available permissions",
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name()},
		},
		func() float64 {
			return float64(entry.Metrics().AvailableConcurrentCalls())
		},
	)
	collectors := []prometheus.Collector{
		maxAllowedConcurrentCallsGauge,
		availableConcurrentCallsGauge,
	}
	return buildRegisterFn(collectors...), buildUnregisterFn(collectors...)
}
