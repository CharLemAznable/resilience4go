package decorator

import (
	. "github.com/CharLemAznable/gogo/fn"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/cache"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/fallback"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
)

type DecorateFunction[T any, R any] interface {
	WithBulkhead(bulkhead.Bulkhead) DecorateFunction[T, R]
	WhenFull(func(T) (R, error)) DecorateFunction[T, R]
	WithTimeLimiter(timelimiter.TimeLimiter) DecorateFunction[T, R]
	WhenTimeout(func(T) (R, error)) DecorateFunction[T, R]
	WithRateLimiter(ratelimiter.RateLimiter) DecorateFunction[T, R]
	WhenOverRate(func(T) (R, error)) DecorateFunction[T, R]
	WithCircuitBreaker(circuitbreaker.CircuitBreaker) DecorateFunction[T, R]
	WhenOverLoad(func(T) (R, error)) DecorateFunction[T, R]
	WithRetry(retry.Retry) DecorateFunction[T, R]
	WhenMaxRetries(func(T) (R, error)) DecorateFunction[T, R]
	WithFallback(func(T) (R, error), func(T, R, error, any) bool) DecorateFunction[T, R]
	WithCache(cache.Cache[T, R]) DecorateFunction[T, R]
	Decorate() Function[T, R]
}

func OfFunction[T any, R any](function Function[T, R]) DecorateFunction[T, R] {
	return &decorateFunction[T, R]{function}
}

type decorateFunction[T any, R any] struct {
	Function[T, R]
}

func (fn *decorateFunction[T, R]) WithBulkhead(entry bulkhead.Bulkhead) DecorateFunction[T, R] {
	return fn.function(bulkhead.DecorateFunction(entry, fn.Function))
}

func (fn *decorateFunction[T, R]) WhenFull(fallbackFn func(T) (R, error)) DecorateFunction[T, R] {
	return fn.function(fallback.DecorateFunctionByType[T, R, *bulkhead.FullError](fn.Function, fallbackFn))
}

func (fn *decorateFunction[T, R]) WithTimeLimiter(entry timelimiter.TimeLimiter) DecorateFunction[T, R] {
	return fn.function(timelimiter.DecorateFunction(entry, fn.Function))
}

func (fn *decorateFunction[T, R]) WhenTimeout(fallbackFn func(T) (R, error)) DecorateFunction[T, R] {
	return fn.function(fallback.DecorateFunctionByType[T, R, *timelimiter.TimeoutError](fn.Function, fallbackFn))
}

func (fn *decorateFunction[T, R]) WithRateLimiter(entry ratelimiter.RateLimiter) DecorateFunction[T, R] {
	return fn.function(ratelimiter.DecorateFunction(entry, fn.Function))
}

func (fn *decorateFunction[T, R]) WhenOverRate(fallbackFn func(T) (R, error)) DecorateFunction[T, R] {
	return fn.function(fallback.DecorateFunctionByType[T, R, *ratelimiter.NotPermittedError](fn.Function, fallbackFn))
}

func (fn *decorateFunction[T, R]) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) DecorateFunction[T, R] {
	return fn.function(circuitbreaker.DecorateFunction(entry, fn.Function))
}

func (fn *decorateFunction[T, R]) WhenOverLoad(fallbackFn func(T) (R, error)) DecorateFunction[T, R] {
	return fn.function(fallback.DecorateFunctionByType[T, R, *circuitbreaker.NotPermittedError](fn.Function, fallbackFn))
}

func (fn *decorateFunction[T, R]) WithRetry(entry retry.Retry) DecorateFunction[T, R] {
	return fn.function(retry.DecorateFunction(entry, fn.Function))
}

func (fn *decorateFunction[T, R]) WhenMaxRetries(fallbackFn func(T) (R, error)) DecorateFunction[T, R] {
	return fn.function(fallback.DecorateFunctionByType[T, R, *retry.MaxRetriesExceeded](fn.Function, fallbackFn))
}

func (fn *decorateFunction[T, R]) WithFallback(
	fallbackFn func(T) (R, error), predicate func(T, R, error, any) bool) DecorateFunction[T, R] {
	return fn.function(fallback.DecorateFunction(fn.Function,
		func(ctx fallback.Context[T, R, error]) (R, error) { return fallbackFn(ctx.Param()) },
		func(ctx fallback.Context[T, R, error]) (bool, fallback.Context[T, R, error]) {
			return predicate(ctx.Param(), ctx.Ret(), ctx.Err(), ctx.Panic()), ctx
		}))
}

func (fn *decorateFunction[T, R]) WithCache(entry cache.Cache[T, R]) DecorateFunction[T, R] {
	return fn.function(cache.DecorateFunction(entry, fn.Function))
}

func (fn *decorateFunction[T, R]) function(function Function[T, R]) DecorateFunction[T, R] {
	fn.Function = function
	return fn
}

func (fn *decorateFunction[T, R]) Decorate() Function[T, R] {
	return fn.Function
}
