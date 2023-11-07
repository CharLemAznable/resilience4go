package promhelper_test

import (
	"github.com/CharLemAznable/resilience4go/promhelper"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestRateLimiterCollectors(t *testing.T) {
	entry := ratelimiter.NewRateLimiter("test") // Create a new ratelimiter entry for testing
	collectors := promhelper.RateLimiterCollectors(entry)
	if len(collectors) != 2 {
		t.Errorf("Expected 2 collectors, but got %d", len(collectors))
	}

	// Assert the first collector
	numberOfWaitingThreadsCollector := collectors[0].(prometheus.GaugeFunc)
	expectedDesc := `Desc{fqName: "resilience4go_ratelimiter_waiting_threads", help: "The number of waiting threads", constLabels: {name="test"}, variableLabels: {}}`
	if numberOfWaitingThreadsCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, numberOfWaitingThreadsCollector.Desc().String())
	}
	m := &dto.Metric{}
	_ = numberOfWaitingThreadsCollector.Write(m)
	expected := &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Gauge: &dto.Gauge{
			Value: proto.Float64(0),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}

	// Assert the second collector
	availablePermissionsCollector := collectors[1].(prometheus.GaugeFunc)
	expectedDesc = `Desc{fqName: "resilience4go_ratelimiter_available_permissions", help: "The number of available permissions", constLabels: {name="test"}, variableLabels: {}}`
	if availablePermissionsCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, availablePermissionsCollector.Desc().String())
	}
	m = &dto.Metric{}
	_ = availablePermissionsCollector.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Gauge: &dto.Gauge{
			Value: proto.Float64(50),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}
}
