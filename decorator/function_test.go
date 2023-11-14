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
		WhenFull(func(_, _ string, fullError *bulkhead.FullError) (string, error) {
			return "", nil
		}).
		WithTimeLimiter(timelimiter.NewTimeLimiter("test")).
		WhenTimeout(func(_, _ string, timeoutError *timelimiter.TimeoutError) (string, error) {
			return "", nil
		}).
		WithRateLimiter(ratelimiter.NewRateLimiter("test")).
		WhenOverRate(func(_, _ string, notPermittedError *ratelimiter.NotPermittedError) (string, error) {
			return "", nil
		}).
		WithCircuitBreaker(circuitbreaker.NewCircuitBreaker("test")).
		WhenOverLoad(func(_, _ string, notPermittedError *circuitbreaker.NotPermittedError) (string, error) {
			return "", nil
		}).
		WithRetry(retry.NewRetry("test")).
		WhenMaxRetries(func(_, _ string, exceeded *retry.MaxRetriesExceeded) (string, error) {
			return "", nil
		}).
		WithFallback(func(_, _ string, err error) (string, error) { return "fallback", nil }).
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
