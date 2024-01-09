package decorator_test

import (
	"errors"
	. "github.com/CharLemAznable/gogo/fn"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/decorator"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
	"testing"
)

func TestDecorateSupplier(t *testing.T) {
	decorateSupplier := decorator.
		OfSupplier(SupplierCast(func() (string, error) {
			return "", errors.New("error")
		})).
		WithBulkhead(bulkhead.NewBulkhead("test")).
		WhenFull(func() (string, error) {
			return "", nil
		}).
		WithTimeLimiter(timelimiter.NewTimeLimiter("test")).
		WhenTimeout(func() (string, error) {
			return "", nil
		}).
		WithRateLimiter(ratelimiter.NewRateLimiter("test")).
		WhenOverRate(func() (string, error) {
			return "", nil
		}).
		WithCircuitBreaker(circuitbreaker.NewCircuitBreaker("test")).
		WhenOverLoad(func() (string, error) {
			return "", nil
		}).
		WithRetry(retry.NewRetry("test")).
		WhenMaxRetries(func() (string, error) {
			return "", nil
		}).
		WithFallback(
			func() (string, error) {
				return "fallback", nil
			},
			func(_ string, err error, _ any) bool {
				return err != nil
			})
	if decorateSupplier == nil {
		t.Error("Expected non-nil decoratedSupplier")
	}
	decoratedSupplier := decorateSupplier.Decorate()
	if decoratedSupplier == nil {
		t.Error("Expected non-nil decoratedSupplier")
	}
	ret, err := decoratedSupplier.CheckedGet()
	if ret != "fallback" {
		t.Errorf("Expected ret is 'fallback', but got '%v'", ret)
	}
	if err != nil {
		t.Errorf("Expected error is nil, but got '%v'", err)
	}
}
