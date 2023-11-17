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

func OfFunction[T any, R any](function function.Function[T, R]) *DecorateFunction[T, R] {
	return &DecorateFunction[T, R]{function}
}

type DecorateFunction[T any, R any] struct {
	function.Function[T, R]
}

func (function *DecorateFunction[T, R]) WithBulkhead(entry bulkhead.Bulkhead) *DecorateFunction[T, R] {
	return OfFunction(bulkhead.DecorateFunction(entry, function.Function))
}

func (function *DecorateFunction[T, R]) WhenFull(fn func(T) (R, error)) *DecorateFunction[T, R] {
	return OfFunction(fallback.DecorateFunctionByType[T, R, *bulkhead.FullError](function.Function, fn))
}

func (function *DecorateFunction[T, R]) WithTimeLimiter(entry timelimiter.TimeLimiter) *DecorateFunction[T, R] {
	return OfFunction(timelimiter.DecorateFunction(entry, function.Function))
}

func (function *DecorateFunction[T, R]) WhenTimeout(fn func(T) (R, error)) *DecorateFunction[T, R] {
	return OfFunction(fallback.DecorateFunctionByType[T, R, *timelimiter.TimeoutError](function.Function, fn))
}

func (function *DecorateFunction[T, R]) WithRateLimiter(entry ratelimiter.RateLimiter) *DecorateFunction[T, R] {
	return OfFunction(ratelimiter.DecorateFunction(entry, function.Function))
}

func (function *DecorateFunction[T, R]) WhenOverRate(fn func(T) (R, error)) *DecorateFunction[T, R] {
	return OfFunction(fallback.DecorateFunctionByType[T, R, *ratelimiter.NotPermittedError](function.Function, fn))
}

func (function *DecorateFunction[T, R]) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) *DecorateFunction[T, R] {
	return OfFunction(circuitbreaker.DecorateFunction(entry, function.Function))
}

func (function *DecorateFunction[T, R]) WhenOverLoad(fn func(T) (R, error)) *DecorateFunction[T, R] {
	return OfFunction(fallback.DecorateFunctionByType[T, R, *circuitbreaker.NotPermittedError](function.Function, fn))
}

func (function *DecorateFunction[T, R]) WithRetry(entry retry.Retry) *DecorateFunction[T, R] {
	return OfFunction(retry.DecorateFunction(entry, function.Function))
}

func (function *DecorateFunction[T, R]) WhenMaxRetries(fn func(T) (R, error)) *DecorateFunction[T, R] {
	return OfFunction(fallback.DecorateFunctionByType[T, R, *retry.MaxRetriesExceeded](function.Function, fn))
}

func (function *DecorateFunction[T, R]) WithFallback(
	fn func(T) (R, error), predicate func(T, R, error, any) bool) *DecorateFunction[T, R] {
	return OfFunction(fallback.DecorateFunction(function.Function,
		func(ctx fallback.Context[T, R, error]) (R, error) { return fn(ctx.Param()) },
		func(ctx fallback.Context[T, R, error]) (bool, fallback.Context[T, R, error]) {
			return predicate(ctx.Param(), ctx.Ret(), ctx.Err(), ctx.Panic()), ctx
		}))
}

func (function *DecorateFunction[T, R]) WithCache(entry cache.Cache[T, R]) *DecorateFunction[T, R] {
	return OfFunction(cache.DecorateFunction(entry, function.Function))
}

func (function *DecorateFunction[T, R]) Decorate() function.Function[T, R] {
	return function.Function
}
