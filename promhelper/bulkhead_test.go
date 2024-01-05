package promhelper_test

import (
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/promhelper"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestBulkheadRegistry(t *testing.T) {
	entry := bulkhead.NewBulkhead("test") // Create a new bulkhead entry for testing
	registerer := &testRegisterer{
		testingT: t,
		testCases: []*metricTestCase{
			{
				name: "TestMaxAllowedConcurrentCalls",
				desc: `Desc{fqName: "resilience4go_bulkhead_max_allowed_concurrent_calls", help: "The maximum number of available permissions", constLabels: {name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Gauge: &dto.Gauge{
						Value: proto.Float64(25),
					},
				},
			},
			{
				name: "TestAvailableConcurrentCalls",
				desc: `Desc{fqName: "resilience4go_bulkhead_available_concurrent_calls", help: "The number of available permissions", constLabels: {name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Gauge: &dto.Gauge{
						Value: proto.Float64(25),
					},
				},
			},
		},
	}
	registerFn, unregisterFn := promhelper.BulkheadRegistry(entry)
	_ = registerFn(registerer)
	unregisterFn(registerer)

	reg := prometheus.NewRegistry()
	if err := registerFn(reg); err != nil {
		t.Errorf("expected none error, but got %v", err)
	}
	unregisterFn(reg)
}
