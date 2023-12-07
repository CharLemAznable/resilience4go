package circuitbreaker

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

type State string

const (
	Closed     State = "CLOSED"
	Open       State = "OPEN"
	HalfOpen   State = "HALF_OPEN"
	Disabled   State = "DISABLED"
	ForcedOpen State = "FORCED_OPEN"
)

type state struct {
	name         State
	allowPublish bool
	attempts     int64
	metrics      *metrics

	acquirePermission func() error
	onError           func(time.Duration)
	onSuccess         func(time.Duration)
	preTransitionHook func()
}

func closed(breaker CircuitBreaker) *state {
	s := &state{
		name:         Closed,
		allowPublish: true,
		attempts:     0,
		metrics:      forClosed(breaker.config()),
	}
	var isClosed atomic.Bool
	isClosed.Store(true)
	checkIfThresholdsExceeded := func(result metricsResult) {
		if exceededThresholds(result) && isClosed.CompareAndSwap(true, false) {
			breaker.publishThresholdsExceededEvent(result, s.metrics)
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

func open(attempts int64, metrics *metrics, breaker CircuitBreaker) *state {
	config := breaker.config()
	s := &state{
		name:         Open,
		allowPublish: true,
		attempts:     attempts,
		metrics:      metrics,
	}
	var isOpen atomic.Bool
	isOpen.Store(true)
	toHalfOpen := func() {
		if isOpen.CompareAndSwap(true, false) {
			_ = breaker.TransitionToHalfOpenState()
		}
	}
	waitDuration := config.waitIntervalFunctionInOpenStateFn(attempts)
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
		timeout, cancelFunc := context.WithTimeout(
			context.Background(), waitDuration)
		go func() {
			select {
			case <-timeout.Done():
				if timeout.Err() == context.DeadlineExceeded {
					toHalfOpen()
				}
			}
		}()
		s.preTransitionHook = cancelFunc
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
		if permittedNumberOfCalls.Add(-1) >= 0 {
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
		if result == belowThresholds {
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
		timeout, cancelFunc := context.WithTimeout(
			context.Background(), config.maxWaitDurationInHalfOpenState)
		go func() {
			select {
			case <-timeout.Done():
				if timeout.Err() == context.DeadlineExceeded {
					toOpen()
				}
			}
		}()
		s.preTransitionHook = cancelFunc
	}
	return s
}

func atomicHalfOpen(breaker CircuitBreaker) (func(), func()) {
	var isHalfOpen atomic.Bool
	isHalfOpen.Store(true)
	toOpen := func() {
		if isHalfOpen.CompareAndSwap(true, false) {
			_ = breaker.TransitionToOpenState()
		}
	}
	toClosed := func() {
		if isHalfOpen.CompareAndSwap(true, false) {
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

func checkStateTransition(name string, fromState, toState State) error {
	if fromState == Closed && toState == HalfOpen {
		return fmt.Errorf("CircuitBreaker '%s' tried an illegal state transition from %s to %s",
			name, fromState, toState)
	}
	return nil
}
