package promhelper

import (
	"github.com/CharLemAznable/resilience4go/cache"
	"github.com/prometheus/client_golang/prometheus"
)

func CacheRegistry[K any, V any](entry cache.Cache[K, V]) (RegisterFn, UnregisterFn) {
	cacheHitsGauge := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "resilience4go_cache_hits",
			Help:        "The number of cache was found",
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name()},
		},
		func() float64 {
			return float64(entry.Metrics().NumberOfCacheHits())
		},
	)
	cacheMissesGauge := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "resilience4go_cache_misses",
			Help:        "The number of cache was not found",
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name()},
		},
		func() float64 {
			return float64(entry.Metrics().NumberOfCacheMisses())
		},
	)
	collectors := []prometheus.Collector{cacheHitsGauge, cacheMissesGauge}
	return buildRegisterFn(collectors...), buildUnregisterFn(collectors...)
}
