package decorator_test

import (
	"errors"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/decorator"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
	"testing"
)

func TestDecorateConsumer(t *testing.T) {
	decoratedConsumer := decorator.
		OfConsumer(func(str string) error {
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
		WithFallback(func(err error) error { return nil }).
		Decorate()

	if decoratedConsumer == nil {
		t.Error("Expected non-nil decoratedConsumer")
	}
	err := decoratedConsumer("test")
	if err != nil {
		t.Errorf("Expected error is nil, but got '%v'", err)
	}
}
