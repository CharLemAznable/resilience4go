package promhelper_test

import (
	"github.com/CharLemAznable/resilience4go/promhelper"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestRateLimiterRegistry(t *testing.T) {
	entry := ratelimiter.NewRateLimiter("test") // Create a new ratelimiter entry for testing
	registerer := &testRegisterer{
		testingT: t,
		testCases: []*metricTestCase{
			{
				name: "TestNumberOfWaitingThreads",
				desc: `Desc{fqName: "resilience4go_ratelimiter_waiting_threads", help: "The number of waiting threads", constLabels: {name="test"}, variableLabels: {}}`,
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
				name: "TestAvailablePermissions",
				desc: `Desc{fqName: "resilience4go_ratelimiter_available_permissions", help: "The number of available permissions", constLabels: {name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Gauge: &dto.Gauge{
						Value: proto.Float64(50),
					},
				},
			},
		},
	}
	registerFn, unregisterFn := promhelper.RateLimiterRegistry(entry)
	_ = registerFn(registerer)
	unregisterFn(registerer)
}
