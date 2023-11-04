package retry_test

import (
	"errors"
	"fmt"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func TestSuccess(t *testing.T) {
	successRetry := retry.NewRetry("success",
		retry.WithMaxAttempts(2),
		retry.WithRecordResultPredicate(nil))
	assert.Equal(t, "success", successRetry.Name())
	listener := successRetry.EventListener()
	listener.OnSuccess(func(event retry.Event) {
		assert.Equal(t, retry.SUCCESS, event.EventType())
		assert.Equal(t, fmt.Sprintf(
			"%v: Retry '%s' recorded a successful retry attempt."+
				" Number of retry attempts: '%d', Last result was: ('%v', '%v').",
			event.CreationTime(), event.RetryName(),
			event.NumOfAttempts(), event.RetVal(), event.RetErr()),
			fmt.Sprintf("%v", event))
	})
	listener.OnError(func(event retry.Event) {
		assert.Fail(t, "should not listen error event")
	})
	listener.OnRetry(func(event retry.Event) {
		assert.Equal(t, retry.RETRY, event.EventType())
		assert.Equal(t, fmt.Sprintf(
			"%v: Retry '%s', waiting %v until attempt '%d'."+
				" Last result was: ('%v', '%v').",
			event.CreationTime(), event.RetryName(), retry.DefaultWaitDuration,
			event.NumOfAttempts(), event.RetVal(), event.RetErr()),
			fmt.Sprintf("%v", event))
	})

	var count atomic.Int64
	fn := func() error {
		if count.Add(1) >= 2 {
			return nil
		}
		return errors.New("error")
	}
	decoratedFn := retry.DecorateRunnable(successRetry, fn)

	err := decoratedFn()
	assert.NoError(t, err)

	time.Sleep(time.Second)
}

func TestError(t *testing.T) {
	successRetry := retry.NewRetry("error",
		retry.WithMaxAttempts(2),
		retry.WithWaitIntervalFunction(nil))
	assert.Equal(t, "error", successRetry.Name())
	listener := successRetry.EventListener()
	listener.OnSuccess(func(event retry.Event) {
		assert.Fail(t, "should not listen success event")
	})
	listener.OnError(func(event retry.Event) {
		assert.Equal(t, retry.ERROR, event.EventType())
		assert.Equal(t, fmt.Sprintf(
			"%v: Retry '%s' recorded a failed retry attempt."+
				" Number of retry attempts: '%d'. Giving up. Last result was: ('%v', '%v').",
			event.CreationTime(), event.RetryName(),
			event.NumOfAttempts(), event.RetVal(), event.RetErr()),
			fmt.Sprintf("%v", event))
	})
	listener.OnRetry(func(event retry.Event) {
		assert.Equal(t, retry.RETRY, event.EventType())
		assert.Equal(t, fmt.Sprintf(
			"%v: Retry '%s', waiting %v until attempt '%d'."+
				" Last result was: ('%v', '%v').",
			event.CreationTime(), event.RetryName(), retry.DefaultWaitDuration,
			event.NumOfAttempts(), event.RetVal(), event.RetErr()),
			fmt.Sprintf("%v", event))
	})

	var count atomic.Int64
	fn := func() error {
		if count.Add(1) >= 3 {
			return nil
		}
		return errors.New("error")
	}
	decoratedFn := retry.DecorateRunnable(successRetry, fn)

	err := decoratedFn()
	assert.Error(t, err)

	time.Sleep(time.Second)
}
