package decorator

import (
	"github.com/CharLemAznable/gofn/runnable"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/fallback"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
)

func OfRunnable(fn func() error) *DecorateRunnable {
	return &DecorateRunnable{fn}
}

type DecorateRunnable struct {
	fn func() error
}

func (runnable *DecorateRunnable) WithBulkhead(entry bulkhead.Bulkhead) *DecorateRunnable {
	return runnable.setFn(bulkhead.DecorateRunnable(entry, runnable.fn))
}

func (runnable *DecorateRunnable) WhenFull(fn func() error) *DecorateRunnable {
	return runnable.setFn(fallback.DecorateRunnableByType[*bulkhead.FullError](runnable.fn, fn))
}

func (runnable *DecorateRunnable) WithTimeLimiter(entry timelimiter.TimeLimiter) *DecorateRunnable {
	return runnable.setFn(timelimiter.DecorateRunnable(entry, runnable.fn))
}

func (runnable *DecorateRunnable) WhenTimeout(fn func() error) *DecorateRunnable {
	return runnable.setFn(fallback.DecorateRunnableByType[*timelimiter.TimeoutError](runnable.fn, fn))
}

func (runnable *DecorateRunnable) WithRateLimiter(entry ratelimiter.RateLimiter) *DecorateRunnable {
	return runnable.setFn(ratelimiter.DecorateRunnable(entry, runnable.fn))
}

func (runnable *DecorateRunnable) WhenOverRate(fn func() error) *DecorateRunnable {
	return runnable.setFn(fallback.DecorateRunnableByType[*ratelimiter.NotPermittedError](runnable.fn, fn))
}

func (runnable *DecorateRunnable) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) *DecorateRunnable {
	return runnable.setFn(circuitbreaker.DecorateRunnable(entry, runnable.fn))
}

func (runnable *DecorateRunnable) WhenOverLoad(fn func() error) *DecorateRunnable {
	return runnable.setFn(fallback.DecorateRunnableByType[*circuitbreaker.NotPermittedError](runnable.fn, fn))
}

func (runnable *DecorateRunnable) WithRetry(entry retry.Retry) *DecorateRunnable {
	return runnable.setFn(retry.DecorateRunnable(entry, runnable.fn))
}

func (runnable *DecorateRunnable) WhenMaxRetries(fn func() error) *DecorateRunnable {
	return runnable.setFn(fallback.DecorateRunnableByType[*retry.MaxRetriesExceeded](runnable.fn, fn))
}

func (runnable *DecorateRunnable) WithFallback(
	fn func() error, predicate func(error, any) bool) *DecorateRunnable {
	return runnable.setFn(fallback.DecorateRunnable(runnable.fn,
		func(ctx fallback.Context[any, any, error]) error { return fn() },
		func(ctx fallback.Context[any, any, error]) (bool, fallback.Context[any, any, error]) {
			return predicate(ctx.Err(), ctx.Panic()), ctx
		}))
}

func (runnable *DecorateRunnable) Decorate() runnable.Runnable {
	return runnable.fn
}

func (runnable *DecorateRunnable) setFn(fn func() error) *DecorateRunnable {
	runnable.fn = fn
	return runnable
}
