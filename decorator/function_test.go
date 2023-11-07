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

func TestDecorateFunction(t *testing.T) {
	decoratedFunction := decorator.
		OfFunction(func(str string) (string, error) {
			return "", errors.New("error")
		}).
		WithBulkhead(bulkhead.NewBulkhead("test")).
		WithTimeLimiter(timelimiter.NewTimeLimiter("test")).
		WithRateLimiter(ratelimiter.NewRateLimiter("test")).
		WithCircuitBreaker(circuitbreaker.NewCircuitBreaker("test")).
		WithRetry(retry.NewRetry("test")).
		WithFallback(func(err error) (string, error) { return "fallback", nil }).
		Decorate()

	if decoratedFunction == nil {
		t.Error("Expected non-nil decoratedFunction")
	}
	ret, err := decoratedFunction("test")
	if ret != "fallback" {
		t.Errorf("Expected ret is 'fallback', but got '%v'", ret)
	}
	if err != nil {
		t.Errorf("Expected error is nil, but got '%v'", err)
	}
}
