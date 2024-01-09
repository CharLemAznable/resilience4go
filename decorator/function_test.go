package decorator_test

import (
	"errors"
	. "github.com/CharLemAznable/gogo/fn"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/cache"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/decorator"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
	"testing"
)

func TestDecorateFunction(t *testing.T) {
	decorateFunction := decorator.
		OfFunction(FunctionCast(func(str string) (string, error) {
			return "", errors.New("error")
		})).
		WithBulkhead(bulkhead.NewBulkhead("test")).
		WhenFull(func(_ string) (string, error) {
			return "", nil
		}).
		WithTimeLimiter(timelimiter.NewTimeLimiter("test")).
		WhenTimeout(func(_ string) (string, error) {
			return "", nil
		}).
		WithRateLimiter(ratelimiter.NewRateLimiter("test")).
		WhenOverRate(func(_ string) (string, error) {
			return "", nil
		}).
		WithCircuitBreaker(circuitbreaker.NewCircuitBreaker("test")).
		WhenOverLoad(func(_ string) (string, error) {
			return "", nil
		}).
		WithRetry(retry.NewRetry("test")).
		WhenMaxRetries(func(_ string) (string, error) {
			return "", nil
		}).
		WithFallback(
			func(_ string) (string, error) {
				return "fallback", nil
			},
			func(_ string, _ string, err error, _ any) bool {
				return err != nil
			}).
		WithCache(cache.NewCache[string, string]("test"))
	if decorateFunction == nil {
		t.Error("Expected non-nil decoratedFunction")
	}
	decoratedFunction := decorateFunction.Decorate()
	if decoratedFunction == nil {
		t.Error("Expected non-nil decoratedFunction")
	}
	ret, err := decoratedFunction.CheckedApply("test")
	if ret != "fallback" {
		t.Errorf("Expected ret is 'fallback', but got '%v'", ret)
	}
	if err != nil {
		t.Errorf("Expected error is nil, but got '%v'", err)
	}
}
