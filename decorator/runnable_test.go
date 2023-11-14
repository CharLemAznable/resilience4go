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
	decoratedRunnable := decorator.
		OfRunnable(func() error {
			return errors.New("error")
		}).
		WithBulkhead(bulkhead.NewBulkhead("test")).
		WhenFull(func(fullError *bulkhead.FullError) error {
			return nil
		}).
		WithTimeLimiter(timelimiter.NewTimeLimiter("test")).
		WhenTimeout(func(timeoutError *timelimiter.TimeoutError) error {
			return nil
		}).
		WithRateLimiter(ratelimiter.NewRateLimiter("test")).
		WhenOverRate(func(notPermittedError *ratelimiter.NotPermittedError) error {
			return nil
		}).
		WithCircuitBreaker(circuitbreaker.NewCircuitBreaker("test")).
		WhenOverLoad(func(notPermittedError *circuitbreaker.NotPermittedError) error {
			return nil
		}).
		WithRetry(retry.NewRetry("test")).
		WhenMaxRetries(func(exceeded *retry.MaxRetriesExceeded) error {
			return nil
		}).
		WithFallback(
			func(err error) error {
				return nil
			},
			func(err error, panic any) (bool, error) {
				return err != nil, err
			}).
		Decorate()

	if decoratedRunnable == nil {
		t.Error("Expected non-nil decoratedRunnable")
	}
	err := decoratedRunnable()
	if err != nil {
		t.Errorf("Expected error is nil, but got '%v'", err)
	}
}
