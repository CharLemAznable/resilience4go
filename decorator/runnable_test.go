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
		WithTimeLimiter(timelimiter.NewTimeLimiter("test")).
		WithRateLimiter(ratelimiter.NewRateLimiter("test")).
		WithCircuitBreaker(circuitbreaker.NewCircuitBreaker("test")).
		WithRetry(retry.NewRetry("test")).
		WithFallback(func(err error) error { return nil }).
		Decorate()

	if decoratedRunnable == nil {
		t.Error("Expected non-nil decoratedRunnable")
	}
	err := decoratedRunnable()
	if err != nil {
		t.Errorf("Expected error is nil, but got '%v'", err)
	}
}
