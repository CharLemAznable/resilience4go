package promhelper

import (
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/prometheus/client_golang/prometheus"
)

func BulkheadCollectors(entry bulkhead.Bulkhead) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "resilience4go_bulkhead_max_allowed_concurrent_calls",
				Help:        "The maximum number of available permissions",
				ConstLabels: prometheus.Labels{labelKeyName: entry.Name()},
			},
			func() float64 {
				return float64(entry.Metrics().MaxAllowedConcurrentCalls())
			},
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "resilience4go_bulkhead_available_concurrent_calls",
				Help:        "The number of available permissions",
				ConstLabels: prometheus.Labels{labelKeyName: entry.Name()},
			},
			func() float64 {
				return float64(entry.Metrics().AvailableConcurrentCalls())
			},
		),
	}
}
