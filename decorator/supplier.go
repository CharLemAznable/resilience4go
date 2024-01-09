package decorator

import (
	. "github.com/CharLemAznable/gogo/fn"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/fallback"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
)

type DecorateSupplier[T any] interface {
	WithBulkhead(bulkhead.Bulkhead) DecorateSupplier[T]
	WhenFull(func() (T, error)) DecorateSupplier[T]
	WithTimeLimiter(timelimiter.TimeLimiter) DecorateSupplier[T]
	WhenTimeout(func() (T, error)) DecorateSupplier[T]
	WithRateLimiter(ratelimiter.RateLimiter) DecorateSupplier[T]
	WhenOverRate(func() (T, error)) DecorateSupplier[T]
	WithCircuitBreaker(circuitbreaker.CircuitBreaker) DecorateSupplier[T]
	WhenOverLoad(func() (T, error)) DecorateSupplier[T]
	WithRetry(retry.Retry) DecorateSupplier[T]
	WhenMaxRetries(func() (T, error)) DecorateSupplier[T]
	WithFallback(func() (T, error), func(T, error, any) bool) DecorateSupplier[T]
	Decorate() Supplier[T]
}

func OfSupplier[T any](supplier Supplier[T]) DecorateSupplier[T] {
	return &decorateSupplier[T]{supplier}
}

type decorateSupplier[T any] struct {
	Supplier[T]
}

func (fn *decorateSupplier[T]) WithBulkhead(entry bulkhead.Bulkhead) DecorateSupplier[T] {
	return fn.supplier(bulkhead.DecorateSupplier(entry, fn.Supplier))
}

func (fn *decorateSupplier[T]) WhenFull(fallbackFn func() (T, error)) DecorateSupplier[T] {
	return fn.supplier(fallback.DecorateSupplierByType[T, *bulkhead.FullError](fn.Supplier, fallbackFn))
}

func (fn *decorateSupplier[T]) WithTimeLimiter(entry timelimiter.TimeLimiter) DecorateSupplier[T] {
	return fn.supplier(timelimiter.DecorateSupplier(entry, fn.Supplier))
}

func (fn *decorateSupplier[T]) WhenTimeout(fallbackFn func() (T, error)) DecorateSupplier[T] {
	return fn.supplier(fallback.DecorateSupplierByType[T, *timelimiter.TimeoutError](fn.Supplier, fallbackFn))
}

func (fn *decorateSupplier[T]) WithRateLimiter(entry ratelimiter.RateLimiter) DecorateSupplier[T] {
	return fn.supplier(ratelimiter.DecorateSupplier(entry, fn.Supplier))
}

func (fn *decorateSupplier[T]) WhenOverRate(fallbackFn func() (T, error)) DecorateSupplier[T] {
	return fn.supplier(fallback.DecorateSupplierByType[T, *ratelimiter.NotPermittedError](fn.Supplier, fallbackFn))
}

func (fn *decorateSupplier[T]) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) DecorateSupplier[T] {
	return fn.supplier(circuitbreaker.DecorateSupplier(entry, fn.Supplier))
}

func (fn *decorateSupplier[T]) WhenOverLoad(fallbackFn func() (T, error)) DecorateSupplier[T] {
	return fn.supplier(fallback.DecorateSupplierByType[T, *circuitbreaker.NotPermittedError](fn.Supplier, fallbackFn))
}

func (fn *decorateSupplier[T]) WithRetry(entry retry.Retry) DecorateSupplier[T] {
	return fn.supplier(retry.DecorateSupplier(entry, fn.Supplier))
}

func (fn *decorateSupplier[T]) WhenMaxRetries(fallbackFn func() (T, error)) DecorateSupplier[T] {
	return fn.supplier(fallback.DecorateSupplierByType[T, *retry.MaxRetriesExceeded](fn.Supplier, fallbackFn))
}

func (fn *decorateSupplier[T]) WithFallback(
	fallbackFn func() (T, error), predicate func(T, error, any) bool) DecorateSupplier[T] {
	return fn.supplier(fallback.DecorateSupplier(fn.Supplier,
		func(ctx fallback.Context[any, T, error]) (T, error) { return fallbackFn() },
		func(ctx fallback.Context[any, T, error]) (bool, fallback.Context[any, T, error]) {
			return predicate(ctx.Ret(), ctx.Err(), ctx.Panic()), ctx
		}))
}

func (fn *decorateSupplier[T]) supplier(supplier Supplier[T]) DecorateSupplier[T] {
	fn.Supplier = supplier
	return fn
}

func (fn *decorateSupplier[T]) Decorate() Supplier[T] {
	return fn.Supplier
}
