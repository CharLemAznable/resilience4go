package retry

import (
	"fmt"
	"time"
)

type EventType string

const (
	SUCCESS EventType = "SUCCESS"
	RETRY   EventType = "RETRY"
	ERROR   EventType = "ERROR"
)

type Event interface {
	fmt.Stringer
	RetryName() string
	CreationTime() time.Time
	EventType() EventType
	NumOfAttempts() int
	RetVal() any
	RetErr() error
}

type SuccessEvent interface {
	Event
}

//goland:noinspection GoNameStartsWithPackageName
type RetryEvent interface {
	Event
	WaitDuration() time.Duration
}

type ErrorEvent interface {
	Event
}

func newSuccessEvent(retryName string, numOfAttempts int, ret any, err error) Event {
	return &successEvent{event{
		retryName:     retryName,
		creationTime:  time.Now(),
		numOfAttempts: numOfAttempts,
		ret:           ret,
		err:           err,
	}}
}

func newRetryEvent(retryName string, numOfAttempts int, ret any, err error, waitDuration time.Duration) Event {
	return &retryEvent{event{
		retryName:     retryName,
		creationTime:  time.Now(),
		numOfAttempts: numOfAttempts,
		ret:           ret,
		err:           err,
	}, waitDuration}
}

func newErrorEvent(retryName string, numOfAttempts int, ret any, err error) Event {
	return &errorEvent{event{
		retryName:     retryName,
		creationTime:  time.Now(),
		numOfAttempts: numOfAttempts,
		ret:           ret,
		err:           err,
	}}
}

type event struct {
	retryName     string
	creationTime  time.Time
	numOfAttempts int
	ret           any
	err           error
}

func (e *event) RetryName() string {
	return e.retryName
}

func (e *event) CreationTime() time.Time {
	return e.creationTime
}

func (e *event) NumOfAttempts() int {
	return e.numOfAttempts
}

func (e *event) RetVal() any {
	return e.ret
}

func (e *event) RetErr() error {
	return e.err
}

type successEvent struct {
	event
}

func (e *successEvent) EventType() EventType {
	return SUCCESS
}

func (e *successEvent) String() string {
	return fmt.Sprintf(
		"%v: Retry '%s' recorded a successful retry attempt."+
			" Number of retry attempts: '%d', Last result was: ('%v', '%v').",
		e.creationTime, e.retryName, e.numOfAttempts, e.ret, e.err)
}

type retryEvent struct {
	event
	waitDuration time.Duration
}

func (e *retryEvent) EventType() EventType {
	return RETRY
}

func (e *retryEvent) WaitDuration() time.Duration {
	return e.waitDuration
}

func (e *retryEvent) String() string {
	return fmt.Sprintf(
		"%v: Retry '%s', waiting %v until attempt '%d'."+
			" Last result was: ('%v', '%v').",
		e.creationTime, e.retryName, e.waitDuration, e.numOfAttempts, e.ret, e.err)
}

type errorEvent struct {
	event
}

func (e *errorEvent) EventType() EventType {
	return ERROR
}

func (e *errorEvent) String() string {
	return fmt.Sprintf(
		"%v: Retry '%s' recorded a failed retry attempt."+
			" Number of retry attempts: '%d'. Giving up. Last result was: ('%v', '%v').",
		e.creationTime, e.retryName, e.numOfAttempts, e.ret, e.err)
}
