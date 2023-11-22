package decorator

import (
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/fallback"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
)

type DecorateSupplier[T any] func() (T, error)

func OfSupplier[T any](fn func() (T, error)) DecorateSupplier[T] {
	return fn
}

func (fn DecorateSupplier[T]) WithBulkhead(entry bulkhead.Bulkhead) DecorateSupplier[T] {
	return bulkhead.DecorateSupplier(entry, fn)
}

func (fn DecorateSupplier[T]) WhenFull(fallbackFn func() (T, error)) DecorateSupplier[T] {
	return fallback.DecorateSupplierByType[T, *bulkhead.FullError](fn, fallbackFn)
}

func (fn DecorateSupplier[T]) WithTimeLimiter(entry timelimiter.TimeLimiter) DecorateSupplier[T] {
	return timelimiter.DecorateSupplier(entry, fn)
}

func (fn DecorateSupplier[T]) WhenTimeout(fallbackFn func() (T, error)) DecorateSupplier[T] {
	return fallback.DecorateSupplierByType[T, *timelimiter.TimeoutError](fn, fallbackFn)
}

func (fn DecorateSupplier[T]) WithRateLimiter(entry ratelimiter.RateLimiter) DecorateSupplier[T] {
	return ratelimiter.DecorateSupplier(entry, fn)
}

func (fn DecorateSupplier[T]) WhenOverRate(fallbackFn func() (T, error)) DecorateSupplier[T] {
	return fallback.DecorateSupplierByType[T, *ratelimiter.NotPermittedError](fn, fallbackFn)
}

func (fn DecorateSupplier[T]) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) DecorateSupplier[T] {
	return circuitbreaker.DecorateSupplier(entry, fn)
}

func (fn DecorateSupplier[T]) WhenOverLoad(fallbackFn func() (T, error)) DecorateSupplier[T] {
	return fallback.DecorateSupplierByType[T, *circuitbreaker.NotPermittedError](fn, fallbackFn)
}

func (fn DecorateSupplier[T]) WithRetry(entry retry.Retry) DecorateSupplier[T] {
	return retry.DecorateSupplier(entry, fn)
}

func (fn DecorateSupplier[T]) WhenMaxRetries(fallbackFn func() (T, error)) DecorateSupplier[T] {
	return fallback.DecorateSupplierByType[T, *retry.MaxRetriesExceeded](fn, fallbackFn)
}

func (fn DecorateSupplier[T]) WithFallback(
	fallbackFn func() (T, error), predicate func(T, error, any) bool) DecorateSupplier[T] {
	return fallback.DecorateSupplier(fn,
		func(ctx fallback.Context[any, T, error]) (T, error) { return fallbackFn() },
		func(ctx fallback.Context[any, T, error]) (bool, fallback.Context[any, T, error]) {
			return predicate(ctx.Ret(), ctx.Err(), ctx.Panic()), ctx
		})
}

func (fn DecorateSupplier[T]) Decorate() func() (T, error) {
	return fn
}
