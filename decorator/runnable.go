package decorator

import (
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/fallback"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
)

type DecorateRunnable func() error

func OfRunnable(fn func() error) DecorateRunnable {
	return fn
}

func (fn DecorateRunnable) WithBulkhead(entry bulkhead.Bulkhead) DecorateRunnable {
	return bulkhead.DecorateRunnable(entry, fn)
}

func (fn DecorateRunnable) WhenFull(fallbackFn func() error) DecorateRunnable {
	return fallback.DecorateRunnableByType[*bulkhead.FullError](fn, fallbackFn)
}

func (fn DecorateRunnable) WithTimeLimiter(entry timelimiter.TimeLimiter) DecorateRunnable {
	return timelimiter.DecorateRunnable(entry, fn)
}

func (fn DecorateRunnable) WhenTimeout(fallbackFn func() error) DecorateRunnable {
	return fallback.DecorateRunnableByType[*timelimiter.TimeoutError](fn, fallbackFn)
}

func (fn DecorateRunnable) WithRateLimiter(entry ratelimiter.RateLimiter) DecorateRunnable {
	return ratelimiter.DecorateRunnable(entry, fn)
}

func (fn DecorateRunnable) WhenOverRate(fallbackFn func() error) DecorateRunnable {
	return fallback.DecorateRunnableByType[*ratelimiter.NotPermittedError](fn, fallbackFn)
}

func (fn DecorateRunnable) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) DecorateRunnable {
	return circuitbreaker.DecorateRunnable(entry, fn)
}

func (fn DecorateRunnable) WhenOverLoad(fallbackFn func() error) DecorateRunnable {
	return fallback.DecorateRunnableByType[*circuitbreaker.NotPermittedError](fn, fallbackFn)
}

func (fn DecorateRunnable) WithRetry(entry retry.Retry) DecorateRunnable {
	return retry.DecorateRunnable(entry, fn)
}

func (fn DecorateRunnable) WhenMaxRetries(fallbackFn func() error) DecorateRunnable {
	return fallback.DecorateRunnableByType[*retry.MaxRetriesExceeded](fn, fallbackFn)
}

func (fn DecorateRunnable) WithFallback(
	fallbackFn func() error, predicate func(error, any) bool) DecorateRunnable {
	return fallback.DecorateRunnable(fn,
		func(ctx fallback.Context[any, any, error]) error { return fallbackFn() },
		func(ctx fallback.Context[any, any, error]) (bool, fallback.Context[any, any, error]) {
			return predicate(ctx.Err(), ctx.Panic()), ctx
		})
}

func (fn DecorateRunnable) Decorate() func() error {
	return fn
}
