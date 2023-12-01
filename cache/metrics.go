package cache

import "github.com/dgraph-io/ristretto"

type Metrics interface {
	NumberOfCacheHits() uint64
	NumberOfCacheMisses() uint64
}

func newMetrics(ristrettoMetrics *ristretto.Metrics) *metrics {
	return &metrics{
		ristrettoMetrics: ristrettoMetrics,
	}
}

type metrics struct {
	ristrettoMetrics *ristretto.Metrics
}

func (m *metrics) NumberOfCacheHits() uint64 {
	return m.ristrettoMetrics.Hits()
}

func (m *metrics) NumberOfCacheMisses() uint64 {
	return m.ristrettoMetrics.Misses()
}
