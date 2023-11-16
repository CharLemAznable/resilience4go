package timelimiter

import (
	"fmt"
	"time"
)

type EventType string

const (
	SUCCESS EventType = "SUCCESS"
	TIMEOUT EventType = "TIMEOUT"
	PANIC   EventType = "PANIC"
)

type Event interface {
	fmt.Stringer
	TimeLimiterName() string
	CreationTime() time.Time
	EventType() EventType
}

type SuccessEvent interface {
	Event
}

type TimeoutEvent interface {
	Event
}

type PanicEvent interface {
	Event
	Panic() any
}

func newSuccessEvent(timeLimiterName string) Event {
	return &successEvent{event{timeLimiterName: timeLimiterName, creationTime: time.Now()}}
}

func newTimeoutEvent(timeLimiterName string) Event {
	return &timeoutEvent{event{timeLimiterName: timeLimiterName, creationTime: time.Now()}}
}

func newPanicEvent(timeLimiterName string, panic any) Event {
	return &panicEvent{event{timeLimiterName: timeLimiterName, creationTime: time.Now()}, panic}
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

type panicEvent struct {
	event
	panic any
}

func (e *panicEvent) EventType() EventType {
	return PANIC
}

func (e *panicEvent) Panic() any {
	return e.panic
}

func (e *panicEvent) String() string {
	return fmt.Sprintf(
		"%v: TimeLimiter '%s' recorded a failure call with panic: %v.",
		e.creationTime, e.timeLimiterName, e.panic)
}
