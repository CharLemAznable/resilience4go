package promhelper_test

import (
	"github.com/CharLemAznable/resilience4go/promhelper"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestRetryCollectors(t *testing.T) {
	entry := retry.NewRetry("test") // Create a new retry entry for testing
	collectors := promhelper.RetryCollectors(entry)
	if len(collectors) != 4 {
		t.Errorf("Expected 4 collectors, but got %d", len(collectors))
	}

	// Assert the first collector
	successWithoutRetryCollector := collectors[0].(prometheus.CounterFunc)
	expectedDesc := `Desc{fqName: "resilience4go_retry_calls", help: "The number of successful calls without a retry attempt", constLabels: {kind="successful_without_retry",name="test"}, variableLabels: {}}`
	if successWithoutRetryCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, successWithoutRetryCollector.Desc().String())
	}
	m := &dto.Metric{}
	_ = successWithoutRetryCollector.Write(m)
	expected := &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("kind"), Value: proto.String("successful_without_retry")},
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Counter: &dto.Counter{
			Value: proto.Float64(0),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}

	// Assert the second collector
	successWithRetryCollector := collectors[1].(prometheus.CounterFunc)
	expectedDesc = `Desc{fqName: "resilience4go_retry_calls", help: "The number of successful calls after a retry attempt", constLabels: {kind="successful_with_retry",name="test"}, variableLabels: {}}`
	if successWithRetryCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, successWithRetryCollector.Desc().String())
	}
	m = &dto.Metric{}
	_ = successWithRetryCollector.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("kind"), Value: proto.String("successful_with_retry")},
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Counter: &dto.Counter{
			Value: proto.Float64(0),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}

	// Assert the third collector
	failedWithoutRetryCollector := collectors[2].(prometheus.CounterFunc)
	expectedDesc = `Desc{fqName: "resilience4go_retry_calls", help: "The number of failed calls without a retry attempt", constLabels: {kind="failed_without_retry",name="test"}, variableLabels: {}}`
	if failedWithoutRetryCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, failedWithoutRetryCollector.Desc().String())
	}
	m = &dto.Metric{}
	_ = failedWithoutRetryCollector.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("kind"), Value: proto.String("failed_without_retry")},
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Counter: &dto.Counter{
			Value: proto.Float64(0),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}

	// Assert the forth collector
	failedWithRetryCollector := collectors[3].(prometheus.CounterFunc)
	expectedDesc = `Desc{fqName: "resilience4go_retry_calls", help: "The number of failed calls after a retry attempt", constLabels: {kind="failed_with_retry",name="test"}, variableLabels: {}}`
	if failedWithRetryCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, failedWithRetryCollector.Desc().String())
	}
	m = &dto.Metric{}
	_ = failedWithRetryCollector.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("kind"), Value: proto.String("failed_with_retry")},
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
