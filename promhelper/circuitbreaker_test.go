package promhelper_test

import (
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/promhelper"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func TestCircuitBreakerCollectors(t *testing.T) {
	entry := circuitbreaker.NewCircuitBreaker("test") // Create a new circuitbreaker entry for testing
	collectors, onSuccess, onError := promhelper.CircuitBreakerCollectors(entry)
	if len(collectors) != 14 {
		t.Errorf("Expected 2 collectors, but got %d", len(collectors))
	}

	// Assert state gauges
	stateCollector := collectors[0].(prometheus.GaugeFunc)
	expectedDesc := `Desc{fqName: "resilience4go_circuitbreaker_state", help: "The states of the circuit breaker", constLabels: {name="test",state="closed"}, variableLabels: {}}`
	if stateCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, stateCollector.Desc().String())
	}
	m := &dto.Metric{}
	_ = stateCollector.Write(m)
	expected := &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
			{Name: proto.String("state"), Value: proto.String("closed")},
		},
		Gauge: &dto.Gauge{
			Value: proto.Float64(1),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}
	stateCollector = collectors[1].(prometheus.GaugeFunc)
	expectedDesc = `Desc{fqName: "resilience4go_circuitbreaker_state", help: "The states of the circuit breaker", constLabels: {name="test",state="open"}, variableLabels: {}}`
	if stateCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, stateCollector.Desc().String())
	}
	m = &dto.Metric{}
	_ = stateCollector.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
			{Name: proto.String("state"), Value: proto.String("open")},
		},
		Gauge: &dto.Gauge{
			Value: proto.Float64(0),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}
	stateCollector = collectors[2].(prometheus.GaugeFunc)
	expectedDesc = `Desc{fqName: "resilience4go_circuitbreaker_state", help: "The states of the circuit breaker", constLabels: {name="test",state="half_open"}, variableLabels: {}}`
	if stateCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, stateCollector.Desc().String())
	}
	m = &dto.Metric{}
	_ = stateCollector.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
			{Name: proto.String("state"), Value: proto.String("half_open")},
		},
		Gauge: &dto.Gauge{
			Value: proto.Float64(0),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}
	stateCollector = collectors[3].(prometheus.GaugeFunc)
	expectedDesc = `Desc{fqName: "resilience4go_circuitbreaker_state", help: "The states of the circuit breaker", constLabels: {name="test",state="disabled"}, variableLabels: {}}`
	if stateCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, stateCollector.Desc().String())
	}
	m = &dto.Metric{}
	_ = stateCollector.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
			{Name: proto.String("state"), Value: proto.String("disabled")},
		},
		Gauge: &dto.Gauge{
			Value: proto.Float64(0),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}
	stateCollector = collectors[4].(prometheus.GaugeFunc)
	expectedDesc = `Desc{fqName: "resilience4go_circuitbreaker_state", help: "The states of the circuit breaker", constLabels: {name="test",state="forced_open"}, variableLabels: {}}`
	if stateCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, stateCollector.Desc().String())
	}
	m = &dto.Metric{}
	_ = stateCollector.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
			{Name: proto.String("state"), Value: proto.String("forced_open")},
		},
		Gauge: &dto.Gauge{
			Value: proto.Float64(0),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}

	// Assert call gauges
	callCollector := collectors[5].(prometheus.GaugeFunc)
	expectedDesc = `Desc{fqName: "resilience4go_circuitbreaker_buffered_calls", help: "The number of buffered successful calls stored in the ring buffer", constLabels: {kind="successful",name="test"}, variableLabels: {}}`
	if callCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, callCollector.Desc().String())
	}
	m = &dto.Metric{}
	_ = callCollector.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("kind"), Value: proto.String("successful")},
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Gauge: &dto.Gauge{
			Value: proto.Float64(0),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}
	callCollector = collectors[6].(prometheus.GaugeFunc)
	expectedDesc = `Desc{fqName: "resilience4go_circuitbreaker_buffered_calls", help: "The number of buffered failed calls stored in the ring buffer", constLabels: {kind="failed",name="test"}, variableLabels: {}}`
	if callCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, callCollector.Desc().String())
	}
	m = &dto.Metric{}
	_ = callCollector.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("kind"), Value: proto.String("failed")},
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Gauge: &dto.Gauge{
			Value: proto.Float64(0),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}
	callCollector = collectors[7].(prometheus.GaugeFunc)
	expectedDesc = `Desc{fqName: "resilience4go_circuitbreaker_slow_calls", help: "The number of slow successful which were slower than a certain threshold", constLabels: {kind="successful",name="test"}, variableLabels: {}}`
	if callCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, callCollector.Desc().String())
	}
	m = &dto.Metric{}
	_ = callCollector.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("kind"), Value: proto.String("successful")},
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Gauge: &dto.Gauge{
			Value: proto.Float64(0),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}
	callCollector = collectors[8].(prometheus.GaugeFunc)
	expectedDesc = `Desc{fqName: "resilience4go_circuitbreaker_slow_calls", help: "The number of slow failed calls which were slower than a certain threshold", constLabels: {kind="failed",name="test"}, variableLabels: {}}`
	if callCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, callCollector.Desc().String())
	}
	m = &dto.Metric{}
	_ = callCollector.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("kind"), Value: proto.String("failed")},
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Gauge: &dto.Gauge{
			Value: proto.Float64(0),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}
	callCollector = collectors[9].(prometheus.GaugeFunc)
	expectedDesc = `Desc{fqName: "resilience4go_circuitbreaker_failure_rate", help: "The failure rate of the circuit breaker", constLabels: {name="test"}, variableLabels: {}}`
	if callCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, callCollector.Desc().String())
	}
	m = &dto.Metric{}
	_ = callCollector.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Gauge: &dto.Gauge{
			Value: proto.Float64(-1),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}
	callCollector = collectors[10].(prometheus.GaugeFunc)
	expectedDesc = `Desc{fqName: "resilience4go_circuitbreaker_slow_call_rate", help: "The slow call rate of the circuit breaker", constLabels: {name="test"}, variableLabels: {}}`
	if callCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, callCollector.Desc().String())
	}
	m = &dto.Metric{}
	_ = callCollector.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Gauge: &dto.Gauge{
			Value: proto.Float64(-1),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}

	// Assert call histograms
	onSuccess(&MockEvent{circuitbreaker.Success, time.Second})
	onError(&MockEvent{circuitbreaker.Error, time.Second * 2})
	callHistogram := collectors[11].(prometheus.Histogram)
	expectedDesc = `Desc{fqName: "resilience4go_circuitbreaker_calls", help: "Total number of successful calls", constLabels: {kind="successful",name="test"}, variableLabels: {}}`
	if callHistogram.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, callHistogram.Desc().String())
	}
	m = &dto.Metric{}
	_ = callHistogram.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("kind"), Value: proto.String("successful")},
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Histogram: &dto.Histogram{
			SampleCount: proto.Uint64(1),
			SampleSum:   proto.Float64(1),
			Bucket: []*dto.Bucket{
				{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(0.005)},
				{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(0.01)},
				{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(0.025)},
				{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(0.05)},
				{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(0.1)},
				{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(0.25)},
				{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(0.5)},
				{CumulativeCount: proto.Uint64(1), UpperBound: proto.Float64(1)},
				{CumulativeCount: proto.Uint64(1), UpperBound: proto.Float64(2.5)},
				{CumulativeCount: proto.Uint64(1), UpperBound: proto.Float64(5)},
				{CumulativeCount: proto.Uint64(1), UpperBound: proto.Float64(10)},
			},
			CreatedTimestamp: m.Histogram.CreatedTimestamp,
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}
	callHistogram = collectors[12].(prometheus.Histogram)
	expectedDesc = `Desc{fqName: "resilience4go_circuitbreaker_calls", help: "Total number of failed calls", constLabels: {kind="failed",name="test"}, variableLabels: {}}`
	if callHistogram.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, callHistogram.Desc().String())
	}
	m = &dto.Metric{}
	_ = callHistogram.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("kind"), Value: proto.String("failed")},
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Histogram: &dto.Histogram{
			SampleCount: proto.Uint64(1),
			SampleSum:   proto.Float64(2),
			Bucket: []*dto.Bucket{
				{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(0.005)},
				{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(0.01)},
				{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(0.025)},
				{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(0.05)},
				{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(0.1)},
				{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(0.25)},
				{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(0.5)},
				{CumulativeCount: proto.Uint64(0), UpperBound: proto.Float64(1)},
				{CumulativeCount: proto.Uint64(1), UpperBound: proto.Float64(2.5)},
				{CumulativeCount: proto.Uint64(1), UpperBound: proto.Float64(5)},
				{CumulativeCount: proto.Uint64(1), UpperBound: proto.Float64(10)},
			},
			CreatedTimestamp: m.Histogram.CreatedTimestamp,
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}

	// Assert notPermitted Counter
	notPermittedCounter := collectors[13].(prometheus.CounterFunc)
	expectedDesc = `Desc{fqName: "resilience4go_circuitbreaker_not_permitted_calls", help: "Total number of not permitted calls", constLabels: {kind="not_permitted",name="test"}, variableLabels: {}}`
	if notPermittedCounter.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, notPermittedCounter.Desc().String())
	}
	m = &dto.Metric{}
	_ = notPermittedCounter.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("kind"), Value: proto.String("not_permitted")},
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Counter: &dto.Counter{
			Value: proto.Float64(0),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}
}

type MockEvent struct {
	eventType circuitbreaker.EventType
	duration  time.Duration
}

func (e *MockEvent) CircuitBreakerName() string {
	return "circuitBreakerName"
}

func (e *MockEvent) CreationTime() time.Time {
	return time.Now()
}

func (e *MockEvent) EventType() circuitbreaker.EventType {
	return e.eventType
}

func (e *MockEvent) Duration() time.Duration {
	return e.duration
}

func (e *MockEvent) String() string {
	return ""
}
