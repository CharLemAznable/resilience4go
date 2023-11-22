package decorator

import (
	"github.com/CharLemAznable/gofn/supplier"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/fallback"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
)

func OfSupplier[T any](fn func() (T, error)) *DecorateSupplier[T] {
	return &DecorateSupplier[T]{fn}
}

type DecorateSupplier[T any] struct {
	fn func() (T, error)
}

func (supplier *DecorateSupplier[T]) WithBulkhead(entry bulkhead.Bulkhead) *DecorateSupplier[T] {
	return supplier.setFn(bulkhead.DecorateSupplier(entry, supplier.fn))
}

func (supplier *DecorateSupplier[T]) WhenFull(fn func() (T, error)) *DecorateSupplier[T] {
	return supplier.setFn(fallback.DecorateSupplierByType[T, *bulkhead.FullError](supplier.fn, fn))
}

func (supplier *DecorateSupplier[T]) WithTimeLimiter(entry timelimiter.TimeLimiter) *DecorateSupplier[T] {
	return supplier.setFn(timelimiter.DecorateSupplier(entry, supplier.fn))
}

func (supplier *DecorateSupplier[T]) WhenTimeout(fn func() (T, error)) *DecorateSupplier[T] {
	return supplier.setFn(fallback.DecorateSupplierByType[T, *timelimiter.TimeoutError](supplier.fn, fn))
}

func (supplier *DecorateSupplier[T]) WithRateLimiter(entry ratelimiter.RateLimiter) *DecorateSupplier[T] {
	return supplier.setFn(ratelimiter.DecorateSupplier(entry, supplier.fn))
}

func (supplier *DecorateSupplier[T]) WhenOverRate(fn func() (T, error)) *DecorateSupplier[T] {
	return supplier.setFn(fallback.DecorateSupplierByType[T, *ratelimiter.NotPermittedError](supplier.fn, fn))
}

func (supplier *DecorateSupplier[T]) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) *DecorateSupplier[T] {
	return supplier.setFn(circuitbreaker.DecorateSupplier(entry, supplier.fn))
}

func (supplier *DecorateSupplier[T]) WhenOverLoad(fn func() (T, error)) *DecorateSupplier[T] {
	return supplier.setFn(fallback.DecorateSupplierByType[T, *circuitbreaker.NotPermittedError](supplier.fn, fn))
}

func (supplier *DecorateSupplier[T]) WithRetry(entry retry.Retry) *DecorateSupplier[T] {
	return supplier.setFn(retry.DecorateSupplier(entry, supplier.fn))
}

func (supplier *DecorateSupplier[T]) WhenMaxRetries(fn func() (T, error)) *DecorateSupplier[T] {
	return supplier.setFn(fallback.DecorateSupplierByType[T, *retry.MaxRetriesExceeded](supplier.fn, fn))
}

func (supplier *DecorateSupplier[T]) WithFallback(
	fn func() (T, error), predicate func(T, error, any) bool) *DecorateSupplier[T] {
	return supplier.setFn(fallback.DecorateSupplier(supplier.fn,
		func(ctx fallback.Context[any, T, error]) (T, error) { return fn() },
		func(ctx fallback.Context[any, T, error]) (bool, fallback.Context[any, T, error]) {
			return predicate(ctx.Ret(), ctx.Err(), ctx.Panic()), ctx
		}))
}

func (supplier *DecorateSupplier[T]) Decorate() supplier.Supplier[T] {
	return supplier.fn
}

func (supplier *DecorateSupplier[T]) setFn(fn func() (T, error)) *DecorateSupplier[T] {
	supplier.fn = fn
	return supplier
}
