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
		WhenFull(func(_ string) error {
			return nil
		}).
		WithTimeLimiter(timelimiter.NewTimeLimiter("test")).
		WhenTimeout(func(_ string) error {
			return nil
		}).
		WithRateLimiter(ratelimiter.NewRateLimiter("test")).
		WhenOverRate(func(_ string) error {
			return nil
		}).
		WithCircuitBreaker(circuitbreaker.NewCircuitBreaker("test")).
		WhenOverLoad(func(_ string) error {
			return nil
		}).
		WithRetry(retry.NewRetry("test")).
		WhenMaxRetries(func(_ string) error {
			return nil
		}).
		WithFallback(
			func(_ string) error {
				return nil
			},
			func(_ string, err error, _ any) bool {
				return err != nil
			}).
		Decorate()

	if decoratedConsumer == nil {
		t.Error("Expected non-nil decoratedConsumer")
	}
	err := decoratedConsumer("test")
	if err != nil {
		t.Errorf("Expected error is nil, but got '%v'", err)
	}
}
