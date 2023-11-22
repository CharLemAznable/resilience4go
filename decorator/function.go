package decorator

import (
	"github.com/CharLemAznable/gofn/function"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/cache"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/fallback"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
)

func OfFunction[T any, R any](fn func(T) (R, error)) *DecorateFunction[T, R] {
	return &DecorateFunction[T, R]{fn}
}

type DecorateFunction[T any, R any] struct {
	fn func(T) (R, error)
}

func (function *DecorateFunction[T, R]) WithBulkhead(entry bulkhead.Bulkhead) *DecorateFunction[T, R] {
	return function.setFn(bulkhead.DecorateFunction(entry, function.fn))
}

func (function *DecorateFunction[T, R]) WhenFull(fn func(T) (R, error)) *DecorateFunction[T, R] {
	return function.setFn(fallback.DecorateFunctionByType[T, R, *bulkhead.FullError](function.fn, fn))
}

func (function *DecorateFunction[T, R]) WithTimeLimiter(entry timelimiter.TimeLimiter) *DecorateFunction[T, R] {
	return function.setFn(timelimiter.DecorateFunction(entry, function.fn))
}

func (function *DecorateFunction[T, R]) WhenTimeout(fn func(T) (R, error)) *DecorateFunction[T, R] {
	return function.setFn(fallback.DecorateFunctionByType[T, R, *timelimiter.TimeoutError](function.fn, fn))
}

func (function *DecorateFunction[T, R]) WithRateLimiter(entry ratelimiter.RateLimiter) *DecorateFunction[T, R] {
	return function.setFn(ratelimiter.DecorateFunction(entry, function.fn))
}

func (function *DecorateFunction[T, R]) WhenOverRate(fn func(T) (R, error)) *DecorateFunction[T, R] {
	return function.setFn(fallback.DecorateFunctionByType[T, R, *ratelimiter.NotPermittedError](function.fn, fn))
}

func (function *DecorateFunction[T, R]) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) *DecorateFunction[T, R] {
	function.fn = circuitbreaker.DecorateFunction(entry, function.fn)
	return function
}

func (function *DecorateFunction[T, R]) WhenOverLoad(fn func(T) (R, error)) *DecorateFunction[T, R] {
	return function.setFn(fallback.DecorateFunctionByType[T, R, *circuitbreaker.NotPermittedError](function.fn, fn))
}

func (function *DecorateFunction[T, R]) WithRetry(entry retry.Retry) *DecorateFunction[T, R] {
	return function.setFn(retry.DecorateFunction(entry, function.fn))
}

func (function *DecorateFunction[T, R]) WhenMaxRetries(fn func(T) (R, error)) *DecorateFunction[T, R] {
	return function.setFn(fallback.DecorateFunctionByType[T, R, *retry.MaxRetriesExceeded](function.fn, fn))
}

func (function *DecorateFunction[T, R]) WithFallback(
	fn func(T) (R, error), predicate func(T, R, error, any) bool) *DecorateFunction[T, R] {
	return function.setFn(fallback.DecorateFunction(function.fn,
		func(ctx fallback.Context[T, R, error]) (R, error) { return fn(ctx.Param()) },
		func(ctx fallback.Context[T, R, error]) (bool, fallback.Context[T, R, error]) {
			return predicate(ctx.Param(), ctx.Ret(), ctx.Err(), ctx.Panic()), ctx
		}))
}

func (function *DecorateFunction[T, R]) WithCache(entry cache.Cache[T, R]) *DecorateFunction[T, R] {
	return function.setFn(cache.DecorateFunction(entry, function.fn))
}

func (function *DecorateFunction[T, R]) Decorate() function.Function[T, R] {
	return function.fn
}

func (function *DecorateFunction[T, R]) setFn(fn func(T) (R, error)) *DecorateFunction[T, R] {
	function.fn = fn
	return function
}
