package bulkhead

import (
	"fmt"
	"time"
)

type EventType string

const (
	PERMITTED EventType = "PERMITTED"
	REJECTED  EventType = "REJECTED"
	FINISHED  EventType = "FINISHED"
)

type Event interface {
	fmt.Stringer
	BulkheadName() string
	CreationTime() time.Time
	EventType() EventType
}

type PermittedEvent interface {
	Event
}

type RejectedEvent interface {
	Event
}

type FinishedEvent interface {
	Event
}

func newPermittedEvent(bulkheadName string) Event {
	return &permittedEvent{event{bulkheadName: bulkheadName, creationTime: time.Now()}}
}

func newRejectedEvent(bulkheadName string) Event {
	return &rejectedEvent{event{bulkheadName: bulkheadName, creationTime: time.Now()}}
}

func newFinishedEvent(bulkheadName string) Event {
	return &finishedEvent{event{bulkheadName: bulkheadName, creationTime: time.Now()}}
}

type event struct {
	bulkheadName string
	creationTime time.Time
}

func (e *event) BulkheadName() string {
	return e.bulkheadName
}

func (e *event) CreationTime() time.Time {
	return e.creationTime
}

type permittedEvent struct {
	event
}

func (e *permittedEvent) EventType() EventType {
	return PERMITTED
}

func (e *permittedEvent) String() string {
	return fmt.Sprintf(
		"%v: Bulkhead '%s' permitted a call.",
		e.creationTime, e.bulkheadName)
}

type rejectedEvent struct {
	event
}

func (e *rejectedEvent) EventType() EventType {
	return REJECTED
}

func (e *rejectedEvent) String() string {
	return fmt.Sprintf(
		"%v: Bulkhead '%s' rejected a call.",
		e.creationTime, e.bulkheadName)
}

type finishedEvent struct {
	event
}

func (e *finishedEvent) EventType() EventType {
	return FINISHED
}

func (e *finishedEvent) String() string {
	return fmt.Sprintf(
		"%v: Bulkhead '%s' has finished a call.",
		e.creationTime, e.bulkheadName)
}
