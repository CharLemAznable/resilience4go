package promhelper

import (
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

const (
	cbKindSuccessful   = "successful"
	cbKindFailed       = "failed"
	cbKindNotPermitted = "not_permitted"
)

func CircuitBreakerCollectors(entry circuitbreaker.CircuitBreaker) (
	collectors []prometheus.Collector, onSuccess circuitbreaker.EventConsumer, onError circuitbreaker.EventConsumer) {
	var result []prometheus.Collector
	result = append(result, stateGauges(entry)...)
	result = append(result, callGauges(entry)...)
	histograms, onSuccess, onError := callHistograms(entry)
	result = append(result, histograms...)
	entry.EventListener().OnSuccess(onSuccess).OnError(onError)
	return append(result, notPermittedCallsCounter(entry)), onSuccess, onError
}

func stateGauges(entry circuitbreaker.CircuitBreaker) []prometheus.Collector {
	return []prometheus.Collector{
		stateGauge(entry, circuitbreaker.Closed),
		stateGauge(entry, circuitbreaker.Open),
		stateGauge(entry, circuitbreaker.HalfOpen),
		stateGauge(entry, circuitbreaker.Disabled),
		stateGauge(entry, circuitbreaker.ForcedOpen),
	}
}

func stateGauge(entry circuitbreaker.CircuitBreaker, state circuitbreaker.State) prometheus.GaugeFunc {
	return prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "resilience4go_circuitbreaker_state",
			Help:        "The states of the circuit breaker",
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyState: strings.ToLower(string(state))},
		},
		func() float64 {
			return isStateGaugeVal(entry.State(), state)
		},
	)
}

func isStateGaugeVal(current, target circuitbreaker.State) float64 {
	if current == target {
		return 1
	}
	return 0
}

func callGauges(entry circuitbreaker.CircuitBreaker) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "resilience4go_circuitbreaker_buffered_calls",
				Help:        "The number of buffered successful calls stored in the ring buffer",
				ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: cbKindSuccessful},
			},
			func() float64 {
				return float64(entry.Metrics().NumberOfSuccessfulCalls())
			},
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "resilience4go_circuitbreaker_buffered_calls",
				Help:        "The number of buffered failed calls stored in the ring buffer",
				ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: cbKindFailed},
			},
			func() float64 {
				return float64(entry.Metrics().NumberOfFailedCalls())
			},
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "resilience4go_circuitbreaker_slow_calls",
				Help:        "The number of slow successful which were slower than a certain threshold",
				ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: cbKindSuccessful},
			},
			func() float64 {
				return float64(entry.Metrics().NumberOfSlowSuccessfulCalls())
			},
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "resilience4go_circuitbreaker_slow_calls",
				Help:        "The number of slow failed calls which were slower than a certain threshold",
				ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: cbKindFailed},
			},
			func() float64 {
				return float64(entry.Metrics().NumberOfSlowFailedCalls())
			},
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "resilience4go_circuitbreaker_failure_rate",
				Help:        "The failure rate of the circuit breaker",
				ConstLabels: prometheus.Labels{labelKeyName: entry.Name()},
			},
			func() float64 {
				return entry.Metrics().FailureRate()
			},
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "resilience4go_circuitbreaker_slow_call_rate",
				Help:        "The slow call rate of the circuit breaker",
				ConstLabels: prometheus.Labels{labelKeyName: entry.Name()},
			},
			func() float64 {
				return entry.Metrics().SlowCallRate()
			},
		),
	}
}

func callHistograms(entry circuitbreaker.CircuitBreaker) (
	[]prometheus.Collector, circuitbreaker.EventConsumer, circuitbreaker.EventConsumer) {
	successfulCallsHistogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:        "resilience4go_circuitbreaker_calls",
			Help:        "Total number of successful calls",
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: cbKindSuccessful},
		})
	onSuccess := func(event circuitbreaker.Event) {
		successfulCallsHistogram.Observe(event.(circuitbreaker.EventWithDuration).Duration().Seconds())
	}
	failedCallsHistogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:        "resilience4go_circuitbreaker_calls",
			Help:        "Total number of failed calls",
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: cbKindFailed},
		})
	onError := func(event circuitbreaker.Event) {
		failedCallsHistogram.Observe(event.(circuitbreaker.EventWithDuration).Duration().Seconds())
	}
	return []prometheus.Collector{successfulCallsHistogram, failedCallsHistogram}, onSuccess, onError
}

func notPermittedCallsCounter(entry circuitbreaker.CircuitBreaker) prometheus.CounterFunc {
	return prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name:        "resilience4go_circuitbreaker_not_permitted_calls",
			Help:        "Total number of not permitted calls",
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: cbKindNotPermitted},
		},
		func() float64 {
			return float64(entry.Metrics().NumberOfNotPermittedCalls())
		},
	)
}
