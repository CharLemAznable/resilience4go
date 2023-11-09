package promhelper_test

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
	"testing"
)

type metricTestCase struct {
	name   string
	desc   string
	metric *dto.Metric
}

type testRegisterer struct {
	testingT  *testing.T
	index     int
	testCases []*metricTestCase
}

func (registerer *testRegisterer) Register(collector prometheus.Collector) error {
	testCase := registerer.testCases[registerer.index]
	registerer.testingT.Run(testCase.name, func(t *testing.T) {
		metric := collector.(prometheus.Metric)
		if metric.Desc().String() != testCase.desc {
			t.Errorf("Expected collector name '%s', but got %s",
				testCase.desc, metric.Desc().String())
		}
		m := &dto.Metric{}
		_ = metric.Write(m)
		if m.Histogram != nil {
			testCase.metric.Histogram.SampleSum = m.Histogram.SampleSum
			testCase.metric.Histogram.CreatedTimestamp = m.Histogram.CreatedTimestamp
		}
		if !proto.Equal(testCase.metric, m) {
			t.Errorf("expected %q, got %q", testCase.metric, m)
		}
	})
	registerer.index++
	return nil
}

func (registerer *testRegisterer) MustRegister(_ ...prometheus.Collector) {
}

func (registerer *testRegisterer) Unregister(_ prometheus.Collector) bool {
	return true
}
