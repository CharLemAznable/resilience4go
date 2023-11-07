package promhelper_test

import (
	"errors"
	"github.com/CharLemAznable/resilience4go/promhelper"
	"github.com/CharLemAznable/resilience4go/timelimiter"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func TestTimeLimiterCollectors(t *testing.T) {
	entry := timelimiter.NewTimeLimiter("test") // Create a new timelimiter entry for testing
	collectors := promhelper.TimeLimiterCollectors(entry)
	if len(collectors) != 3 {
		t.Errorf("Expected 3 collectors, but got %d", len(collectors))
	}

	fn := timelimiter.DecorateRunnable(entry, func() error {
		panic("panic error")
	})
	func() {
		defer func() {
			if r := recover(); r != nil {
				// ignored
			}
		}()
		_ = fn()
	}()
	fn = timelimiter.DecorateRunnable(entry, func() error {
		time.Sleep(time.Second * 2)
		return nil
	})
	_ = fn()
	fn = timelimiter.DecorateRunnable(entry, func() error {
		time.Sleep(time.Millisecond * 500)
		return errors.New("error")
	})
	_ = fn()

	// Assert the first collector
	successCountCollector := collectors[0].(prometheus.CounterFunc)
	expectedDesc := `Desc{fqName: "resilience4go_timelimiter_calls", help: "The number of successful calls", constLabels: {kind="successful",name="test"}, variableLabels: {}}`
	if successCountCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, successCountCollector.Desc().String())
	}
	m := &dto.Metric{}
	_ = successCountCollector.Write(m)
	expected := &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("kind"), Value: proto.String("successful")},
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Counter: &dto.Counter{
			Value: proto.Float64(1),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}

	// Assert the second collector
	timeoutCountCollector := collectors[1].(prometheus.CounterFunc)
	expectedDesc = `Desc{fqName: "resilience4go_timelimiter_calls", help: "The number of timed out calls", constLabels: {kind="timeout",name="test"}, variableLabels: {}}`
	if timeoutCountCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, timeoutCountCollector.Desc().String())
	}
	m = &dto.Metric{}
	_ = timeoutCountCollector.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("kind"), Value: proto.String("timeout")},
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Counter: &dto.Counter{
			Value: proto.Float64(1),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}

	// Assert the third collector
	failureCountCollector := collectors[2].(prometheus.CounterFunc)
	expectedDesc = `Desc{fqName: "resilience4go_timelimiter_calls", help: "The number of failed calls", constLabels: {kind="failed",name="test"}, variableLabels: {}}`
	if failureCountCollector.Desc().String() != expectedDesc {
		t.Errorf("Expected collector name '%s', but got %s",
			expectedDesc, failureCountCollector.Desc().String())
	}
	m = &dto.Metric{}
	_ = failureCountCollector.Write(m)
	expected = &dto.Metric{
		Label: []*dto.LabelPair{
			{Name: proto.String("kind"), Value: proto.String("failed")},
			{Name: proto.String("name"), Value: proto.String(entry.Name())},
		},
		Counter: &dto.Counter{
			Value: proto.Float64(1),
		},
	}
	if !proto.Equal(expected, m) {
		t.Errorf("expected %q, got %q", expected, m)
	}
}
