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
		WhenFull(func(_ string, fullError *bulkhead.FullError) error {
			return nil
		}).
		WithTimeLimiter(timelimiter.NewTimeLimiter("test")).
		WhenTimeout(func(_ string, timeoutError *timelimiter.TimeoutError) error {
			return nil
		}).
		WithRateLimiter(ratelimiter.NewRateLimiter("test")).
		WhenOverRate(func(_ string, notPermittedError *ratelimiter.NotPermittedError) error {
			return nil
		}).
		WithCircuitBreaker(circuitbreaker.NewCircuitBreaker("test")).
		WhenOverLoad(func(_ string, notPermittedError *circuitbreaker.NotPermittedError) error {
			return nil
		}).
		WithRetry(retry.NewRetry("test")).
		WhenMaxRetries(func(_ string, exceeded *retry.MaxRetriesExceeded) error {
			return nil
		}).
		WithFallback(func(_ string, err error) error { return nil }).
		Decorate()

	if decoratedConsumer == nil {
		t.Error("Expected non-nil decoratedConsumer")
	}
	err := decoratedConsumer("test")
	if err != nil {
		t.Errorf("Expected error is nil, but got '%v'", err)
	}
}
