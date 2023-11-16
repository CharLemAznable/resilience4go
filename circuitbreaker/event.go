package circuitbreaker

import (
	"fmt"
	"time"
)

type EventType string

const (
	Success              EventType = "SUCCESS"
	Error                EventType = "ERROR"
	NotPermitted         EventType = "NOT_PERMITTED"
	StateTransition      EventType = "STATE_TRANSITION"
	FailureRateExceeded  EventType = "FAILURE_RATE_EXCEEDED"
	SlowCallRateExceeded EventType = "SLOW_CALL_RATE_EXCEEDED"
)

type Event interface {
	fmt.Stringer
	CircuitBreakerName() string
	CreationTime() time.Time
	EventType() EventType
}

type SuccessEvent interface {
	Event
	Duration() time.Duration
}

type ErrorEvent interface {
	Event
	Duration() time.Duration
	Ret() any
	Err() error
}

type NotPermittedEvent interface {
	Event
}

type StateTransitionEvent interface {
	Event
	FromState() State
	ToState() State
}

type FailureRateExceededEvent interface {
	Event
	FailureRate() float64
}

type SlowCallRateExceededEvent interface {
	Event
	SlowCallRate() float64
}

func newSuccessEvent(circuitBreakerName string, duration time.Duration) Event {
	return &successEvent{event{circuitBreakerName: circuitBreakerName, creationTime: time.Now()}, duration}
}

func newErrorEvent(circuitBreakerName string, duration time.Duration, ret any, err error) Event {
	return &errorEvent{event{circuitBreakerName: circuitBreakerName, creationTime: time.Now()}, duration, ret, err}
}

func newNotPermittedEvent(circuitBreakerName string) Event {
	return &notPermittedEvent{event{circuitBreakerName: circuitBreakerName, creationTime: time.Now()}}
}

func newStateTransitionEvent(circuitBreakerName string, fromState, toState State) Event {
	return &stateTransitionEvent{event{circuitBreakerName: circuitBreakerName, creationTime: time.Now()}, fromState, toState}
}

func newFailureRateExceededEvent(circuitBreakerName string, failureRate float64) Event {
	return &failureRateExceededEvent{event{circuitBreakerName: circuitBreakerName, creationTime: time.Now()}, failureRate}
}

func newSlowCallRateExceededEvent(circuitBreakerName string, slowCallRate float64) Event {
	return &slowCallRateExceededEvent{event{circuitBreakerName: circuitBreakerName, creationTime: time.Now()}, slowCallRate}
}

type event struct {
	circuitBreakerName string
	creationTime       time.Time
}

func (e *event) CircuitBreakerName() string {
	return e.circuitBreakerName
}

func (e *event) CreationTime() time.Time {
	return e.creationTime
}

type successEvent struct {
	event
	duration time.Duration
}

func (e *successEvent) EventType() EventType {
	return Success
}

func (e *successEvent) Duration() time.Duration {
	return e.duration
}

func (e *successEvent) String() string {
	return fmt.Sprintf(
		"%v: CircuitBreaker '%s' recorded a successful call. Elapsed time: %v",
		e.creationTime, e.circuitBreakerName, e.duration)
}

type errorEvent struct {
	event
	duration time.Duration
	ret      any
	err      error
}

func (e *errorEvent) EventType() EventType {
	return Error
}

func (e *errorEvent) Duration() time.Duration {
	return e.duration
}

func (e *errorEvent) Ret() any {
	return e.ret
}

func (e *errorEvent) Err() error {
	return e.err
}

func (e *errorEvent) String() string {
	return fmt.Sprintf(
		"%v: CircuitBreaker '%s' recorded an error ret '%v' with error: '%s'. Elapsed time: %v",
		e.creationTime, e.circuitBreakerName, e.ret, e.err.Error(), e.duration)
}

type notPermittedEvent struct {
	event
}

func (e *notPermittedEvent) EventType() EventType {
	return NotPermitted
}

func (e *notPermittedEvent) String() string {
	return fmt.Sprintf(
		"%v: CircuitBreaker '%s' recorded a call which was not permitted.",
		e.creationTime, e.circuitBreakerName)
}

type stateTransitionEvent struct {
	event
	fromState, toState State
}

func (e *stateTransitionEvent) EventType() EventType {
	return StateTransition
}

func (e *stateTransitionEvent) FromState() State {
	return e.fromState
}

func (e *stateTransitionEvent) ToState() State {
	return e.toState
}

func (e *stateTransitionEvent) String() string {
	return fmt.Sprintf(
		"%v: CircuitBreaker '%s' changed state from %s to %s",
		e.creationTime, e.circuitBreakerName, e.fromState, e.toState)
}

type failureRateExceededEvent struct {
	event
	failureRate float64
}

func (e *failureRateExceededEvent) EventType() EventType {
	return FailureRateExceeded
}

func (e *failureRateExceededEvent) FailureRate() float64 {
	return e.failureRate
}

func (e *failureRateExceededEvent) String() string {
	return fmt.Sprintf(
		"%v: CircuitBreaker '%s' exceeded failure rate threshold. Current failure rate: %f",
		e.creationTime, e.circuitBreakerName, e.failureRate)
}

type slowCallRateExceededEvent struct {
	event
	slowCallRate float64
}

func (e *slowCallRateExceededEvent) EventType() EventType {
	return SlowCallRateExceeded
}

func (e *slowCallRateExceededEvent) SlowCallRate() float64 {
	return e.slowCallRate
}

func (e *slowCallRateExceededEvent) String() string {
	return fmt.Sprintf(
		"%v: CircuitBreaker '%s' exceeded slow call rate threshold. Current slow call rate: %f",
		e.creationTime, e.circuitBreakerName, e.slowCallRate)
}
