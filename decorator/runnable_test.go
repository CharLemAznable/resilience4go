package decorator_test

import (
	"errors"
	"github.com/CharLemAznable/resilience4go/decorator"
	"testing"

	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
)

func TestDecorateRunnable(t *testing.T) {
	decorateRunnable := decorator.
		OfRunnable(func() error {
			return errors.New("error")
		}).
		WithBulkhead(bulkhead.NewBulkhead("test")).
		WhenFull(func() error {
			return nil
		}).
		WithTimeLimiter(timelimiter.NewTimeLimiter("test")).
		WhenTimeout(func() error {
			return nil
		}).
		WithRateLimiter(ratelimiter.NewRateLimiter("test")).
		WhenOverRate(func() error {
			return nil
		}).
		WithCircuitBreaker(circuitbreaker.NewCircuitBreaker("test")).
		WhenOverLoad(func() error {
			return nil
		}).
		WithRetry(retry.NewRetry("test")).
		WhenMaxRetries(func() error {
			return nil
		}).
		WithFallback(
			func() error {
				return nil
			},
			func(err error, _ any) bool {
				return err != nil
			})
	if decorateRunnable == nil {
		t.Error("Expected non-nil decoratedRunnable")
	}
	err := decorateRunnable()
	if err != nil {
		t.Errorf("Expected error is nil, but got '%v'", err)
	}

	decoratedRunnable := decorateRunnable.Decorate()
	if decoratedRunnable == nil {
		t.Error("Expected non-nil decoratedRunnable")
	}
	err = decoratedRunnable()
	if err != nil {
		t.Errorf("Expected error is nil, but got '%v'", err)
	}
}
