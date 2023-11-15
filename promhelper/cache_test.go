package promhelper_test

import (
	"github.com/CharLemAznable/resilience4go/cache"
	"github.com/CharLemAznable/resilience4go/promhelper"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestCacheRegistry(t *testing.T) {
	entry := cache.NewCache[string, int]("test") // Create a new cache entry for testing
	registerer := &testRegisterer{
		testingT: t,
		testCases: []*metricTestCase{
			{
				name: "TestNumberOfCacheHits",
				desc: `Desc{fqName: "resilience4go_cache_hits", help: "The number of cache was found", constLabels: {name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Gauge: &dto.Gauge{
						Value: proto.Float64(0),
					},
				},
			},
			{
				name: "TestNumberOfCacheMisses",
				desc: `Desc{fqName: "resilience4go_cache_misses", help: "The number of cache was not found", constLabels: {name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Gauge: &dto.Gauge{
						Value: proto.Float64(0),
					},
				},
			},
		},
	}
	registerFn, unregisterFn := promhelper.CacheRegistry(entry)
	_ = registerFn(registerer)
	unregisterFn(registerer)
}
