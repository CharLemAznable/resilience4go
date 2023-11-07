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
		WithTimeLimiter(timelimiter.NewTimeLimiter("test")).
		WithRateLimiter(ratelimiter.NewRateLimiter("test")).
		WithCircuitBreaker(circuitbreaker.NewCircuitBreaker("test")).
		WithRetry(retry.NewRetry("test")).
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
