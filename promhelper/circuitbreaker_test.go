package promhelper_test

import (
	"errors"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/promhelper"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func TestCircuitBreakerRegistry(t *testing.T) {
	entry := circuitbreaker.NewCircuitBreaker("test") // Create a new circuitbreaker entry for testing
	registerer := &testRegisterer{
		testingT: t,
		testCases: []*metricTestCase{
			{
				name: "TestStateClosed",
				desc: `Desc{fqName: "resilience4go_circuitbreaker_state", help: "The states of the circuit breaker", constLabels: {name="test",state="closed"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
						{Name: proto.String("state"), Value: proto.String("closed")},
					},
					Gauge: &dto.Gauge{
						Value: proto.Float64(1),
					},
				},
			},
			{
				name: "TestStateOpen",
				desc: `Desc{fqName: "resilience4go_circuitbreaker_state", help: "The states of the circuit breaker", constLabels: {name="test",state="open"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
						{Name: proto.String("state"), Value: proto.String("open")},
					},
					Gauge: &dto.Gauge{
						Value: proto.Float64(0),
					},
				},
			},
			{
				name: "TestStateHalfOpen",
				desc: `Desc{fqName: "resilience4go_circuitbreaker_state", help: "The states of the circuit breaker", constLabels: {name="test",state="half_open"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
						{Name: proto.String("state"), Value: proto.String("half_open")},
					},
					Gauge: &dto.Gauge{
						Value: proto.Float64(0),
					},
				},
			},
			{
				name: "TestStateDisabled",
				desc: `Desc{fqName: "resilience4go_circuitbreaker_state", help: "The states of the circuit breaker", constLabels: {name="test",state="disabled"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
						{Name: proto.String("state"), Value: proto.String("disabled")},
					},
					Gauge: &dto.Gauge{
						Value: proto.Float64(0),
					},
				},
			},
			{
				name: "TestStateForcedOpen",
				desc: `Desc{fqName: "resilience4go_circuitbreaker_state", help: "The states of the circuit breaker", constLabels: {name="test",state="forced_open"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
						{Name: proto.String("state"), Value: proto.String("forced_open")},
					},
					Gauge: &dto.Gauge{
						Value: proto.Float64(0),
					},
				},
			},
			{
				name: "TestNumberOfSuccessfulCalls",
				desc: `Desc{fqName: "resilience4go_circuitbreaker_buffered_calls_successful", help: "The number of buffered successful calls stored in the ring buffer", constLabels: {kind="successful",name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("kind"), Value: proto.String("successful")},
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Gauge: &dto.Gauge{
						Value: proto.Float64(0),
					},
				},
			},
			{
				name: "TestNumberOfFailedCalls",
				desc: `Desc{fqName: "resilience4go_circuitbreaker_buffered_calls_failed", help: "The number of buffered failed calls stored in the ring buffer", constLabels: {kind="failed",name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("kind"), Value: proto.String("failed")},
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Gauge: &dto.Gauge{
						Value: proto.Float64(0),
					},
				},
			},
			{
				name: "TestNumberOfSlowSuccessfulCalls",
				desc: `Desc{fqName: "resilience4go_circuitbreaker_slow_calls_successful", help: "The number of slow successful which were slower than a certain threshold", constLabels: {kind="successful",name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("kind"), Value: proto.String("successful")},
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Gauge: &dto.Gauge{
						Value: proto.Float64(0),
					},
				},
			},
			{
				name: "TestNumberOfSlowFailedCalls",
				desc: `Desc{fqName: "resilience4go_circuitbreaker_slow_calls_failed", help: "The number of slow failed calls which were slower than a certain threshold", constLabels: {kind="failed",name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("kind"), Value: proto.String("failed")},
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Gauge: &dto.Gauge{
						Value: proto.Float64(0),
					},
				},
			},
			{
				name: "TestFailureRate",
				desc: `Desc{fqName: "resilience4go_circuitbreaker_failure_rate", help: "The failure rate of the circuit breaker", constLabels: {name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Gauge: &dto.Gauge{
						Value: proto.Float64(-1),
					},
				},
			},
			{
				name: "TestSlowCallRate",
				desc: `Desc{fqName: "resilience4go_circuitbreaker_slow_call_rate", help: "The slow call rate of the circuit breaker", constLabels: {name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Gauge: &dto.Gauge{
						Value: proto.Float64(-1),
					},
				},
			},
			{
				name: "TestSuccessfulCallsHistogram",
				desc: `Desc{fqName: "resilience4go_circuitbreaker_calls_successful", help: "Total number of successful calls", constLabels: {kind="successful",name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("kind"), Value: proto.String("successful")},
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Histogram: &dto.Histogram{
						SampleCount: proto.Uint64(0),
						Bucket: []*dto.Bucket{
							{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(float64(time.Second))},
							{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(float64(time.Second * 5))},
							{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(float64(time.Second * 10))},
						},
					},
				},
			},
			{
				name: "TestFailedCallsHistogram",
				desc: `Desc{fqName: "resilience4go_circuitbreaker_calls_failed", help: "Total number of failed calls", constLabels: {kind="failed",name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("kind"), Value: proto.String("failed")},
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Histogram: &dto.Histogram{
						SampleCount: proto.Uint64(0),
						Bucket: []*dto.Bucket{
							{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(float64(time.Second))},
							{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(float64(time.Second * 5))},
							{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(float64(time.Second * 10))},
						},
					},
				},
			},
			{
				name: "TestNumberOfNotPermittedCalls",
				desc: `Desc{fqName: "resilience4go_circuitbreaker_not_permitted_calls", help: "Total number of not permitted calls", constLabels: {kind="not_permitted",name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("kind"), Value: proto.String("not_permitted")},
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Counter: &dto.Counter{
						Value: proto.Float64(0),
					},
				},
			},
		},
	}
	registerFn, unregisterFn := promhelper.CircuitBreakerRegistry(entry,
		float64(time.Second), float64(time.Second*5), float64(time.Second*10))
	_ = registerFn(registerer)

	_ = circuitbreaker.DecorateRunnable(entry, func() error {
		time.Sleep(time.Second * 2)
		return nil
	})()
	_ = circuitbreaker.DecorateRunnable(entry, func() error {
		time.Sleep(time.Second * 6)
		return errors.New("error")
	})()
	time.Sleep(time.Second)
	registerer.index = 0
	registerer.testCases[5].metric.Gauge = &dto.Gauge{
		Value: proto.Float64(1),
	}
	registerer.testCases[6].metric.Gauge = &dto.Gauge{
		Value: proto.Float64(1),
	}
	registerer.testCases[11].metric.Histogram = &dto.Histogram{
		SampleCount: proto.Uint64(1),
		Bucket: []*dto.Bucket{
			{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(float64(time.Second))},
			{CumulativeCount: proto.Uint64(1), UpperBound: proto.Float64(float64(time.Second * 5))},
			{CumulativeCount: proto.Uint64(1), UpperBound: proto.Float64(float64(time.Second * 10))},
		},
	}
	registerer.testCases[12].metric.Histogram = &dto.Histogram{
		SampleCount: proto.Uint64(1),
		Bucket: []*dto.Bucket{
			{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(float64(time.Second))},
			{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(float64(time.Second * 5))},
			{CumulativeCount: proto.Uint64(1), UpperBound: proto.Float64(float64(time.Second * 10))},
		},
	}
	_ = registerFn(registerer)

	unregisterFn(registerer)
}
