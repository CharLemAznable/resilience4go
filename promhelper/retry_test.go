package promhelper_test

import (
	"github.com/CharLemAznable/resilience4go/promhelper"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestRetryRegistry(t *testing.T) {
	entry := retry.NewRetry("test") // Create a new retry entry for testing
	registerer := &testRegisterer{
		testingT: t,
		testCases: []*metricTestCase{
			{
				name: "TestNumberOfSuccessfulCallsWithoutRetryAttempt",
				desc: `Desc{fqName: "resilience4go_retry_calls", help: "The number of successful/failed calls with/without retry", constLabels: {kind="successful_without_retry",name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("kind"), Value: proto.String("successful_without_retry")},
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Counter: &dto.Counter{
						Value: proto.Float64(0),
					},
				},
			},
			{
				name: "TestNumberOfSuccessfulCallsWithRetryAttempt",
				desc: `Desc{fqName: "resilience4go_retry_calls", help: "The number of successful/failed calls with/without retry", constLabels: {kind="successful_with_retry",name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("kind"), Value: proto.String("successful_with_retry")},
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Counter: &dto.Counter{
						Value: proto.Float64(0),
					},
				},
			},
			{
				name: "TestNumberOfFailedCallsWithoutRetryAttempt",
				desc: `Desc{fqName: "resilience4go_retry_calls", help: "The number of successful/failed calls with/without retry", constLabels: {kind="failed_without_retry",name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("kind"), Value: proto.String("failed_without_retry")},
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Counter: &dto.Counter{
						Value: proto.Float64(0),
					},
				},
			},
			{
				name: "TestNumberOfFailedCallsWithRetryAttempt",
				desc: `Desc{fqName: "resilience4go_retry_calls", help: "The number of successful/failed calls with/without retry", constLabels: {kind="failed_with_retry",name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("kind"), Value: proto.String("failed_with_retry")},
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Counter: &dto.Counter{
						Value: proto.Float64(0),
					},
				},
			},
		},
	}
	registerFn, unregisterFn := promhelper.RetryRegistry(entry)
	_ = registerFn(registerer)
	unregisterFn(registerer)

	reg := prometheus.NewRegistry()
	if err := registerFn(reg); err != nil {
		t.Errorf("expected none error, but got %v", err)
	}
	unregisterFn(reg)
}
