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

type DecorateRunnable interface {
	WithBulkhead(bulkhead.Bulkhead) DecorateRunnable
	WhenFull(func() error) DecorateRunnable
	WithTimeLimiter(timelimiter.TimeLimiter) DecorateRunnable
	WhenTimeout(func() error) DecorateRunnable
	WithRateLimiter(ratelimiter.RateLimiter) DecorateRunnable
	WhenOverRate(func() error) DecorateRunnable
	WithCircuitBreaker(circuitbreaker.CircuitBreaker) DecorateRunnable
	WhenOverLoad(func() error) DecorateRunnable
	WithRetry(retry.Retry) DecorateRunnable
	WhenMaxRetries(func() error) DecorateRunnable
	WithFallback(func() error, func(error, any) bool) DecorateRunnable
	Decorate() Runnable
}

func OfRunnable(runnable Runnable) DecorateRunnable {
	return &decorateRunnable{runnable}
}

type decorateRunnable struct {
	Runnable
}

func (fn *decorateRunnable) WithBulkhead(entry bulkhead.Bulkhead) DecorateRunnable {
	return fn.runnable(bulkhead.DecorateRunnable(entry, fn.Runnable))
}

func (fn *decorateRunnable) WhenFull(fallbackFn func() error) DecorateRunnable {
	return fn.runnable(fallback.DecorateRunnableByType[*bulkhead.FullError](fn.Runnable, fallbackFn))
}

func (fn *decorateRunnable) WithTimeLimiter(entry timelimiter.TimeLimiter) DecorateRunnable {
	return fn.runnable(timelimiter.DecorateRunnable(entry, fn.Runnable))
}

func (fn *decorateRunnable) WhenTimeout(fallbackFn func() error) DecorateRunnable {
	return fn.runnable(fallback.DecorateRunnableByType[*timelimiter.TimeoutError](fn.Runnable, fallbackFn))
}

func (fn *decorateRunnable) WithRateLimiter(entry ratelimiter.RateLimiter) DecorateRunnable {
	return fn.runnable(ratelimiter.DecorateRunnable(entry, fn.Runnable))
}

func (fn *decorateRunnable) WhenOverRate(fallbackFn func() error) DecorateRunnable {
	return fn.runnable(fallback.DecorateRunnableByType[*ratelimiter.NotPermittedError](fn.Runnable, fallbackFn))
}

func (fn *decorateRunnable) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) DecorateRunnable {
	return fn.runnable(circuitbreaker.DecorateRunnable(entry, fn.Runnable))
}

func (fn *decorateRunnable) WhenOverLoad(fallbackFn func() error) DecorateRunnable {
	return fn.runnable(fallback.DecorateRunnableByType[*circuitbreaker.NotPermittedError](fn.Runnable, fallbackFn))
}

func (fn *decorateRunnable) WithRetry(entry retry.Retry) DecorateRunnable {
	return fn.runnable(retry.DecorateRunnable(entry, fn.Runnable))
}

func (fn *decorateRunnable) WhenMaxRetries(fallbackFn func() error) DecorateRunnable {
	return fn.runnable(fallback.DecorateRunnableByType[*retry.MaxRetriesExceeded](fn.Runnable, fallbackFn))
}

func (fn *decorateRunnable) WithFallback(
	fallbackFn func() error, predicate func(error, any) bool) DecorateRunnable {
	return fn.runnable(fallback.DecorateRunnable(fn.Runnable,
		func(ctx fallback.Context[any, any, error]) error { return fallbackFn() },
		func(ctx fallback.Context[any, any, error]) (bool, fallback.Context[any, any, error]) {
			return predicate(ctx.Err(), ctx.Panic()), ctx
		}))
}

func (fn *decorateRunnable) runnable(runnable Runnable) DecorateRunnable {
	fn.Runnable = runnable
	return fn
}

func (fn *decorateRunnable) Decorate() Runnable {
	return fn.Runnable
}
