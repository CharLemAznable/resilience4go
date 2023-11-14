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

func TestDecorateSupplier(t *testing.T) {
	decoratedSupplier := decorator.
		OfSupplier(func() (string, error) {
			return "", errors.New("error")
		}).
		WithBulkhead(bulkhead.NewBulkhead("test")).
		WhenFull(func(_ string, fullError *bulkhead.FullError) (string, error) {
			return "", nil
		}).
		WithTimeLimiter(timelimiter.NewTimeLimiter("test")).
		WhenTimeout(func(_ string, timeoutError *timelimiter.TimeoutError) (string, error) {
			return "", nil
		}).
		WithRateLimiter(ratelimiter.NewRateLimiter("test")).
		WhenOverRate(func(_ string, notPermittedError *ratelimiter.NotPermittedError) (string, error) {
			return "", nil
		}).
		WithCircuitBreaker(circuitbreaker.NewCircuitBreaker("test")).
		WhenOverLoad(func(_ string, notPermittedError *circuitbreaker.NotPermittedError) (string, error) {
			return "", nil
		}).
		WithRetry(retry.NewRetry("test")).
		WhenMaxRetries(func(_ string, exceeded *retry.MaxRetriesExceeded) (string, error) {
			return "", nil
		}).
		WithFallback(func(_ string, err error) (string, error) { return "fallback", nil }).
		Decorate()

	if decoratedSupplier == nil {
		t.Error("Expected non-nil decoratedSupplier")
	}
	ret, err := decoratedSupplier()
	if ret != "fallback" {
		t.Errorf("Expected ret is 'fallback', but got '%v'", ret)
	}
	if err != nil {
		t.Errorf("Expected error is nil, but got '%v'", err)
	}
}
