package decorator

import (
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/cache"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/fallback"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
)

type DecorateFunction[T any, R any] func(T) (R, error)

func OfFunction[T any, R any](fn func(T) (R, error)) DecorateFunction[T, R] {
	return fn
}

func (fn DecorateFunction[T, R]) WithBulkhead(entry bulkhead.Bulkhead) DecorateFunction[T, R] {
	return bulkhead.DecorateFunction(entry, fn)
}

func (fn DecorateFunction[T, R]) WhenFull(fallbackFn func(T) (R, error)) DecorateFunction[T, R] {
	return fallback.DecorateFunctionByType[T, R, *bulkhead.FullError](fn, fallbackFn)
}

func (fn DecorateFunction[T, R]) WithTimeLimiter(entry timelimiter.TimeLimiter) DecorateFunction[T, R] {
	return timelimiter.DecorateFunction(entry, fn)
}

func (fn DecorateFunction[T, R]) WhenTimeout(fallbackFn func(T) (R, error)) DecorateFunction[T, R] {
	return fallback.DecorateFunctionByType[T, R, *timelimiter.TimeoutError](fn, fallbackFn)
}

func (fn DecorateFunction[T, R]) WithRateLimiter(entry ratelimiter.RateLimiter) DecorateFunction[T, R] {
	return ratelimiter.DecorateFunction(entry, fn)
}

func (fn DecorateFunction[T, R]) WhenOverRate(fallbackFn func(T) (R, error)) DecorateFunction[T, R] {
	return fallback.DecorateFunctionByType[T, R, *ratelimiter.NotPermittedError](fn, fallbackFn)
}

func (fn DecorateFunction[T, R]) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) DecorateFunction[T, R] {
	return circuitbreaker.DecorateFunction(entry, fn)
}

func (fn DecorateFunction[T, R]) WhenOverLoad(fallbackFn func(T) (R, error)) DecorateFunction[T, R] {
	return fallback.DecorateFunctionByType[T, R, *circuitbreaker.NotPermittedError](fn, fallbackFn)
}

func (fn DecorateFunction[T, R]) WithRetry(entry retry.Retry) DecorateFunction[T, R] {
	return retry.DecorateFunction(entry, fn)
}

func (fn DecorateFunction[T, R]) WhenMaxRetries(fallbackFn func(T) (R, error)) DecorateFunction[T, R] {
	return fallback.DecorateFunctionByType[T, R, *retry.MaxRetriesExceeded](fn, fallbackFn)
}

func (fn DecorateFunction[T, R]) WithFallback(
	fallbackFn func(T) (R, error), predicate func(T, R, error, any) bool) DecorateFunction[T, R] {
	return fallback.DecorateFunction(fn,
		func(ctx fallback.Context[T, R, error]) (R, error) { return fallbackFn(ctx.Param()) },
		func(ctx fallback.Context[T, R, error]) (bool, fallback.Context[T, R, error]) {
			return predicate(ctx.Param(), ctx.Ret(), ctx.Err(), ctx.Panic()), ctx
		})
}

func (fn DecorateFunction[T, R]) WithCache(entry cache.Cache[T, R]) DecorateFunction[T, R] {
	return cache.DecorateFunction(entry, fn)
}

func (fn DecorateFunction[T, R]) Decorate() func(T) (R, error) {
	return fn
}
