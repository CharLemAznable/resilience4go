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

func CircuitBreakerRegistry(entry circuitbreaker.CircuitBreaker, histogramBuckets ...float64) (RegisterFn, UnregisterFn) {
	var collectors []prometheus.Collector
	collectors = append(collectors, stateGauges(entry)...)
	collectors = append(collectors, callGauges(entry)...)
	histograms, onSuccess, onError := callHistograms(entry, histogramBuckets...)
	collectors = append(collectors, histograms...)
	collectors = append(collectors, notPermittedCallsCounter(entry))
	registerFn, unregisterFn := buildRegisterFn(collectors...), buildUnregisterFn(collectors...)
	return func(registerer prometheus.Registerer) error {
			entry.EventListener().OnSuccess(onSuccess).OnError(onError)
			return registerFn(registerer)
		}, func(registerer prometheus.Registerer) bool {
			entry.EventListener().Dismiss(onSuccess).Dismiss(onError)
			return unregisterFn(registerer)
		}
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

const (
	bufferedCallsName = "resilience4go_circuitbreaker_buffered_calls"
	bufferedCallsHelp = "The number of buffered calls stored in the ring buffer"
	slowCallsName     = "resilience4go_circuitbreaker_slow_calls"
	slowCallsHelp     = "The number of slow calls which were slower than a certain threshold"
)

func callGauges(entry circuitbreaker.CircuitBreaker) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        bufferedCallsName,
				Help:        bufferedCallsHelp,
				ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: cbKindSuccessful},
			},
			func() float64 {
				return float64(entry.Metrics().NumberOfSuccessfulCalls())
			},
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        bufferedCallsName,
				Help:        bufferedCallsHelp,
				ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: cbKindFailed},
			},
			func() float64 {
				return float64(entry.Metrics().NumberOfFailedCalls())
			},
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        slowCallsName,
				Help:        slowCallsHelp,
				ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: cbKindSuccessful},
			},
			func() float64 {
				return float64(entry.Metrics().NumberOfSlowSuccessfulCalls())
			},
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        slowCallsName,
				Help:        slowCallsHelp,
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

const (
	callsHistogramName = "resilience4go_circuitbreaker_calls"
	callsHistogramHelp = "Total number of successful/failed calls"
)

func callHistograms(entry circuitbreaker.CircuitBreaker, histogramBuckets ...float64) (
	[]prometheus.Collector, func(circuitbreaker.SuccessEvent), func(circuitbreaker.ErrorEvent)) {
	buckets := prometheus.DefBuckets
	if len(histogramBuckets) > 0 {
		buckets = histogramBuckets
	}
	successfulCallsHistogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:        callsHistogramName,
			Help:        callsHistogramHelp,
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: cbKindSuccessful},
			Buckets:     buckets,
		})
	onSuccess := func(event circuitbreaker.SuccessEvent) {
		successfulCallsHistogram.Observe(float64(event.Duration()))
	}
	failedCallsHistogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:        callsHistogramName,
			Help:        callsHistogramHelp,
			ConstLabels: prometheus.Labels{labelKeyName: entry.Name(), labelKeyKind: cbKindFailed},
			Buckets:     buckets,
		})
	onError := func(event circuitbreaker.ErrorEvent) {
		failedCallsHistogram.Observe(float64(event.Duration()))
	}
	return []prometheus.Collector{successfulCallsHistogram, failedCallsHistogram}, onSuccess, onError
}

func notPermittedCallsCounter(entry circuitbreaker.CircuitBreaker) prometheus.Collector {
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
