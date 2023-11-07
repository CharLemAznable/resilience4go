package retry_test

import (
	"errors"
	"fmt"
	"github.com/CharLemAznable/resilience4go/retry"
	"sync/atomic"
	"testing"
	"time"
)

func TestSuccess(t *testing.T) {
	successRetry := retry.NewRetry("success",
		retry.WithMaxAttempts(2),
		retry.WithRecordResultPredicate(nil))
	if successRetry.Name() != "success" {
		t.Errorf("Expected retry name 'success', but got '%s'", successRetry.Name())
	}
	listener := successRetry.EventListener()
	listener.OnSuccess(func(event retry.Event) {
		if event.EventType() != retry.SUCCESS {
			t.Errorf("Expected event type SUCCESS, but got '%s'", event.EventType())
		}
		expected := fmt.Sprintf(
			"%v: Retry '%s' recorded a successful retry attempt."+
				" Number of retry attempts: '%d', Last result was: ('%v', '%v').",
			event.CreationTime(), event.RetryName(),
			event.NumOfAttempts(), event.RetVal(), event.RetErr())
		if fmt.Sprintf("%v", event) != expected {
			t.Errorf("Expected event string '%s', but got '%v'", expected, event)
		}
	})
	listener.OnError(func(event retry.Event) {
		t.Error("Should not listen error event")
	})
	listener.OnRetry(func(event retry.Event) {
		if event.EventType() != retry.RETRY {
			t.Errorf("Expected event type RETRY, but got '%s'", event.EventType())
		}
		expected := fmt.Sprintf(
			"%v: Retry '%s', waiting %v until attempt '%d'."+
				" Last result was: ('%v', '%v').",
			event.CreationTime(), event.RetryName(), retry.DefaultWaitDuration,
			event.NumOfAttempts(), event.RetVal(), event.RetErr())
		if fmt.Sprintf("%v", event) != expected {
			t.Errorf("Expected event string '%s', but got '%v'", expected, event)
		}
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
	if err != nil {
		t.Errorf("Expected nil error, but got '%v'", err)
	}

	time.Sleep(time.Second)
	metrics := successRetry.Metrics()
	if metrics.NumberOfSuccessfulCallsWithoutRetryAttempt() != 0 {
		t.Errorf("Expected successful calls without retry attempt '0', but got '%d'",
			metrics.NumberOfSuccessfulCallsWithoutRetryAttempt())
	}
	if metrics.NumberOfSuccessfulCallsWithRetryAttempt() != 1 {
		t.Errorf("Expected successful calls with retry attempt '1', but got '%d'",
			metrics.NumberOfSuccessfulCallsWithRetryAttempt())
	}
	if metrics.NumberOfFailedCallsWithoutRetryAttempt() != 0 {
		t.Errorf("Expected failed calls without retry attempt '0', but got '%d'",
			metrics.NumberOfFailedCallsWithoutRetryAttempt())
	}
	if metrics.NumberOfFailedCallsWithRetryAttempt() != 0 {
		t.Errorf("Expected failed calls with retry attempt '0', but got '%d'",
			metrics.NumberOfFailedCallsWithRetryAttempt())
	}
}

func TestError(t *testing.T) {
	successRetry := retry.NewRetry("error",
		retry.WithMaxAttempts(2),
		retry.WithWaitIntervalFunction(nil))
	if successRetry.Name() != "error" {
		t.Errorf("Expected retry name 'error', but got '%s'", successRetry.Name())
	}
	listener := successRetry.EventListener()
	listener.OnSuccess(func(event retry.Event) {
		t.Error("Should not listen success event")
	})
	listener.OnError(func(event retry.Event) {
		if event.EventType() != retry.ERROR {
			t.Errorf("Expected event type ERROR, but got '%s'", event.EventType())
		}
		expected := fmt.Sprintf(
			"%v: Retry '%s' recorded a failed retry attempt."+
				" Number of retry attempts: '%d'. Giving up. Last result was: ('%v', '%v').",
			event.CreationTime(), event.RetryName(),
			event.NumOfAttempts(), event.RetVal(), event.RetErr())
		if fmt.Sprintf("%v", event) != expected {
			t.Errorf("Expected event string '%s', but got '%v'", expected, event)
		}
	})
	listener.OnRetry(func(event retry.Event) {
		if event.EventType() != retry.RETRY {
			t.Errorf("Expected event type RETRY, but got '%s'", event.EventType())
		}
		expected := fmt.Sprintf(
			"%v: Retry '%s', waiting %v until attempt '%d'."+
				" Last result was: ('%v', '%v').",
			event.CreationTime(), event.RetryName(), retry.DefaultWaitDuration,
			event.NumOfAttempts(), event.RetVal(), event.RetErr())
		if fmt.Sprintf("%v", event) != expected {
			t.Errorf("Expected event string '%s', but got '%v'", expected, event)
		}
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
	if err == nil {
		t.Error("Expected non-nil error")
	}

	time.Sleep(time.Second)
	metrics := successRetry.Metrics()
	if metrics.NumberOfSuccessfulCallsWithoutRetryAttempt() != 0 {
		t.Errorf("Expected successful calls without retry attempt '0', but got '%d'",
			metrics.NumberOfSuccessfulCallsWithoutRetryAttempt())
	}
	if metrics.NumberOfSuccessfulCallsWithRetryAttempt() != 0 {
		t.Errorf("Expected successful calls with retry attempt '0', but got '%d'",
			metrics.NumberOfSuccessfulCallsWithRetryAttempt())
	}
	if metrics.NumberOfFailedCallsWithoutRetryAttempt() != 0 {
		t.Errorf("Expected failed calls without retry attempt '0', but got '%d'",
			metrics.NumberOfFailedCallsWithoutRetryAttempt())
	}
	if metrics.NumberOfFailedCallsWithRetryAttempt() != 1 {
		t.Errorf("Expected failed calls with retry attempt '1', but got '%d'",
			metrics.NumberOfFailedCallsWithRetryAttempt())
	}
}
