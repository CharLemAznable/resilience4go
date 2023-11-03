package circuitbreaker

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

type stateName string

const (
	Closed     stateName = "CLOSED"
	Open       stateName = "OPEN"
	HalfOpen   stateName = "HALF_OPEN"
	Disabled   stateName = "DISABLED"
	ForcedOpen stateName = "FORCED_OPEN"
)

type state struct {
	name         stateName
	allowPublish bool
	attempts     int64
	metrics      Metrics

	acquirePermission func() error
	onError           func(time.Duration)
	onSuccess         func(time.Duration)
	preTransitionHook func() // nil-able
}

func closed(breaker CircuitBreaker) *state {
	s := &state{
		name:         Closed,
		allowPublish: true,
		attempts:     0,
		metrics:      forClosed(breaker.config()),
	}
	var isClosed atomic.Int32
	isClosed.Store(1)
	checkIfThresholdsExceeded := func(result metricsResult) {
		if exceededThresholds(result) && isClosed.CompareAndSwap(1, 0) {
			_ = breaker.TransitionToOpenState()
		}
	}
	s.onError = func(duration time.Duration) {
		checkIfThresholdsExceeded(s.metrics.onError(duration))
	}
	s.onSuccess = func(duration time.Duration) {
		checkIfThresholdsExceeded(s.metrics.onSuccess(duration))
	}
	return s
}

func open(attempts int64, metrics Metrics, breaker CircuitBreaker) *state {
	config := breaker.config()
	s := &state{
		name:         Open,
		allowPublish: true,
		attempts:     attempts,
		metrics:      metrics,
	}
	var isOpen atomic.Int32
	isOpen.Store(1)
	toHalfOpen := func() {
		if isOpen.CompareAndSwap(1, 0) {
			_ = breaker.TransitionToHalfOpenState()
		}
	}
	waitDuration := config.waitIntervalFunctionInOpenStateFn()(attempts)
	retryAfterWaitDuration := time.Now().Add(waitDuration)
	s.acquirePermission = func() error {
		if time.Now().After(retryAfterWaitDuration) {
			toHalfOpen()
			return breaker.acquirePermission()
		}
		s.metrics.onCallNotPermitted()
		return &NotPermittedError{
			name:      breaker.Name(),
			stateName: Open}
	}
	s.onError = func(duration time.Duration) {
		s.metrics.onError(duration)
	}
	s.onSuccess = func(duration time.Duration) {
		s.metrics.onSuccess(duration)
	}
	if config.automaticTransitionFromOpenToHalfOpenEnabled {
		var done atomic.Int32
		timer := time.NewTimer(waitDuration)
		cancel := make(chan bool, 1)
		go func() {
			select {
			case <-timer.C:
				if done.CompareAndSwap(0, 1) {
					toHalfOpen()
				}
			case <-cancel:
			}
		}()
		s.preTransitionHook = func() {
			if done.CompareAndSwap(0, 1) {
				cancel <- true
			}
		}
	}
	return s
}

func halfOpen(attempts int64, breaker CircuitBreaker) *state {
	config := breaker.config()
	permittedNumber := config.permittedNumberOfCallsInHalfOpenState
	s := &state{
		name:         HalfOpen,
		allowPublish: true,
		attempts:     attempts,
		metrics:      forHalfOpen(permittedNumber, config),
	}
	var permittedNumberOfCalls atomic.Int64
	permittedNumberOfCalls.Store(permittedNumber)
	s.acquirePermission = func() error {
		if getAndUpdateInt64(&permittedNumberOfCalls,
			permittedNumberDecrement) > 0 {
			return nil
		}
		s.metrics.onCallNotPermitted()
		return &NotPermittedError{
			name:      breaker.Name(),
			stateName: HalfOpen}
	}
	toOpen, toClosed := atomicHalfOpen(breaker)
	checkIfThresholdsExceeded := func(result metricsResult) {
		if exceededThresholds(result) {
			toOpen()
		}
		if result == BelowThresholds {
			toClosed()
		}
	}
	s.onError = func(duration time.Duration) {
		checkIfThresholdsExceeded(s.metrics.onError(duration))
	}
	s.onSuccess = func(duration time.Duration) {
		checkIfThresholdsExceeded(s.metrics.onSuccess(duration))
	}
	if config.maxWaitDurationInHalfOpenState > 0 {
		var done atomic.Int32
		timer := time.NewTimer(config.maxWaitDurationInHalfOpenState)
		cancel := make(chan bool, 1)
		go func() {
			select {
			case <-timer.C:
				if done.CompareAndSwap(0, 1) {
					toOpen()
				}
			case <-cancel:
			}
		}()
		s.preTransitionHook = func() {
			if done.CompareAndSwap(0, 1) {
				cancel <- true
			}
		}
	}
	return s
}

func permittedNumberDecrement(current int64) int64 {
	if current == 0 {
		return current
	} else {
		return current - 1
	}
}

func atomicHalfOpen(breaker CircuitBreaker) (func(), func()) {
	var isHalfOpen atomic.Int32
	isHalfOpen.Store(1)
	toOpen := func() {
		if isHalfOpen.CompareAndSwap(1, 0) {
			_ = breaker.TransitionToOpenState()
		}
	}
	toClosed := func() {
		if isHalfOpen.CompareAndSwap(1, 0) {
			_ = breaker.TransitionToClosedState()
		}
	}
	return toOpen, toClosed
}

func disabled(breaker CircuitBreaker) *state {
	return &state{
		name:         Disabled,
		allowPublish: false,
		attempts:     0,
		metrics:      forDisabled(breaker.config()),
	}
}

func forcedOpen(attempts int64, breaker CircuitBreaker) *state {
	s := &state{
		name:         ForcedOpen,
		allowPublish: false,
		attempts:     attempts,
		metrics:      forForcedOpen(breaker.config()),
	}
	s.acquirePermission = func() error {
		s.metrics.onCallNotPermitted()
		return &NotPermittedError{
			name:      breaker.Name(),
			stateName: ForcedOpen}
	}
	return s
}

type stateTransition struct {
	fromState stateName
	toState   stateName
}

func newStateTransition(name string, fromState, toState stateName) (*stateTransition, error) {
	if fromState == Closed && toState == HalfOpen {
		return nil, errors.New(fmt.Sprintf(
			"CircuitBreaker '%s' tried an illegal state transition from %s to %s",
			name, fromState, toState))
	}
	return &stateTransition{
		fromState: fromState,
		toState:   toState,
	}, nil
}
