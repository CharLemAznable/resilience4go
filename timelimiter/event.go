package timelimiter

import (
	"fmt"
	"time"
)

type EventType string

const (
	SUCCESS EventType = "SUCCESS"
	TIMEOUT EventType = "TIMEOUT"
	FAILURE EventType = "FAILURE"
)

type Event interface {
	TimeLimiterName() string
	CreationTime() time.Time
	EventType() EventType
}

func newSuccessEvent(timeLimiterName string) Event {
	return &successEvent{event{timeLimiterName: timeLimiterName, creationTime: time.Now()}}
}

func newTimeoutEvent(timeLimiterName string) Event {
	return &timeoutEvent{event{timeLimiterName: timeLimiterName, creationTime: time.Now()}}
}

func newFailureEvent(timeLimiterName string, error any) Event {
	return &failureEvent{event{timeLimiterName: timeLimiterName, creationTime: time.Now()}, error}
}

type event struct {
	timeLimiterName string
	creationTime    time.Time
}

func (e *event) TimeLimiterName() string {
	return e.timeLimiterName
}

func (e *event) CreationTime() time.Time {
	return e.creationTime
}

type successEvent struct {
	event
}

func (e *successEvent) EventType() EventType {
	return SUCCESS
}

func (e *successEvent) String() string {
	return fmt.Sprintf(
		"%v: TimeLimiter '%s' recorded a successful call.",
		e.creationTime, e.timeLimiterName)
}

type timeoutEvent struct {
	event
}

func (e *timeoutEvent) EventType() EventType {
	return TIMEOUT
}

func (e *timeoutEvent) String() string {
	return fmt.Sprintf(
		"%v: TimeLimiter '%s' recorded a timeout exception.",
		e.creationTime, e.timeLimiterName)
}

type failureEvent struct {
	event
	error any
}

func (e *failureEvent) EventType() EventType {
	return FAILURE
}

func (e *failureEvent) String() string {
	return fmt.Sprintf(
		"%v: TimeLimiter '%s' recorded a failure call with panic: %v.",
		e.creationTime, e.timeLimiterName, e.error)
}
