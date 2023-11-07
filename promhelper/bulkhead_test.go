package promhelper_test

import (
	"testing"

	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/promhelper"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
)

func TestBulkheadCollectors(t *testing.T) {
	entry := bulkhead.NewBulkhead("test") // Create a new bulkhead entry for testing
	collectors := promhelper.BulkheadCollectors(entry)
	if len(collectors) != 2 {
		t.Errorf("Expected 2 collectors, but got %d", len(collectors))
	}

	// Assert the first collector
	maxAllowedConcurrentCallsCollector := collectors[0].(prometheus.GaugeFunc)
	expectedDesc := `Desc{fqName: "resilience4go_bulkhead_max_allowed_concurrent_calls", help: "The maximum number of available permissions", constLabels: {name="test"}, variableLabels: {}}`
	if maxAllowedConcurrentCallsCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, maxAllowedConcurrentCallsCollector.Desc().String())
	}
	m := &dto.Metric{}
	_ = maxAllowedConcurrentCallsCollector.Write(m)
	expected := &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Gauge: &dto.Gauge{
			Value: proto.Float64(25),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}

	// Assert the second collector
	availableConcurrentCallsCollector := collectors[1].(prometheus.GaugeFunc)
	expectedDesc = `Desc{fqName: "resilience4go_bulkhead_available_concurrent_calls", help: "The number of available permissions", constLabels: {name="test"}, variableLabels: {}}`
	if availableConcurrentCallsCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, availableConcurrentCallsCollector.Desc().String())
	}
	m = &dto.Metric{}
	_ = availableConcurrentCallsCollector.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Gauge: &dto.Gauge{
			Value: proto.Float64(25),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}
}
