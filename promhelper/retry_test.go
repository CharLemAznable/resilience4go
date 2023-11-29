package promhelper_test

import (
	"github.com/CharLemAznable/resilience4go/promhelper"
	"github.com/CharLemAznable/resilience4go/retry"
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
				desc: `Desc{fqName: "resilience4go_retry_calls_successful_without_retry", help: "The number of successful calls without a retry attempt", constLabels: {kind="successful_without_retry",name="test"}, variableLabels: {}}`,
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
				desc: `Desc{fqName: "resilience4go_retry_calls_successful_with_retry", help: "The number of successful calls after a retry attempt", constLabels: {kind="successful_with_retry",name="test"}, variableLabels: {}}`,
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
				desc: `Desc{fqName: "resilience4go_retry_calls_failed_without_retry", help: "The number of failed calls without a retry attempt", constLabels: {kind="failed_without_retry",name="test"}, variableLabels: {}}`,
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
				desc: `Desc{fqName: "resilience4go_retry_calls_failed_with_retry", help: "The number of failed calls after a retry attempt", constLabels: {kind="failed_with_retry",name="test"}, variableLabels: {}}`,
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
}
