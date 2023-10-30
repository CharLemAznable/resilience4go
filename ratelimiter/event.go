package ratelimiter

import (
	"fmt"
	"time"
)

type EventType string

const (
	SUCCESSFUL EventType = "SUCCESSFUL_ACQUIRE"
	FAILED     EventType = "FAILED_ACQUIRE"
)

type Event interface {
	RateLimiterName() string
	CreationTime() time.Time
	EventType() EventType
}

func newSuccessEvent(rateLimiterName string) Event {
	return &successEvent{event{rateLimiterName: rateLimiterName, creationTime: time.Now()}}
}

func newFailureEvent(rateLimiterName string) Event {
	return &failureEvent{event{rateLimiterName: rateLimiterName, creationTime: time.Now()}}
}

type event struct {
	rateLimiterName string
	creationTime    time.Time
}

func (e *event) RateLimiterName() string {
	return e.rateLimiterName
}

func (e *event) CreationTime() time.Time {
	return e.creationTime
}

type successEvent struct {
	event
}

func (e *successEvent) EventType() EventType {
	return SUCCESSFUL
}

func (e *successEvent) String() string {
	return fmt.Sprintf(
		"RateLimiterEvent{type=%s, rateLimiterName='%s', creationTime=%v}",
		e.EventType(), e.rateLimiterName, e.creationTime)
}

type failureEvent struct {
	event
}

func (e *failureEvent) EventType() EventType {
	return FAILED
}

func (e *failureEvent) String() string {
	return fmt.Sprintf(
		"RateLimiterEvent{type=%s, rateLimiterName='%s', creationTime=%v}",
		e.EventType(), e.rateLimiterName, e.creationTime)
}
