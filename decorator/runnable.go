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

func OfRunnable(runnable runnable.Runnable) *DecorateRunnable {
	return &DecorateRunnable{runnable}
}

type DecorateRunnable struct {
	runnable.Runnable
}

func (runnable *DecorateRunnable) WithBulkhead(entry bulkhead.Bulkhead) *DecorateRunnable {
	return OfRunnable(bulkhead.DecorateRunnable(entry, runnable.Runnable))
}

func (runnable *DecorateRunnable) WhenFull(fn func() error) *DecorateRunnable {
	return OfRunnable(fallback.DecorateRunnableByType[*bulkhead.FullError](runnable.Runnable, fn))
}

func (runnable *DecorateRunnable) WithTimeLimiter(entry timelimiter.TimeLimiter) *DecorateRunnable {
	return OfRunnable(timelimiter.DecorateRunnable(entry, runnable.Runnable))
}

func (runnable *DecorateRunnable) WhenTimeout(fn func() error) *DecorateRunnable {
	return OfRunnable(fallback.DecorateRunnableByType[*timelimiter.TimeoutError](runnable.Runnable, fn))
}

func (runnable *DecorateRunnable) WithRateLimiter(entry ratelimiter.RateLimiter) *DecorateRunnable {
	return OfRunnable(ratelimiter.DecorateRunnable(entry, runnable.Runnable))
}

func (runnable *DecorateRunnable) WhenOverRate(fn func() error) *DecorateRunnable {
	return OfRunnable(fallback.DecorateRunnableByType[*ratelimiter.NotPermittedError](runnable.Runnable, fn))
}

func (runnable *DecorateRunnable) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) *DecorateRunnable {
	return OfRunnable(circuitbreaker.DecorateRunnable(entry, runnable.Runnable))
}

func (runnable *DecorateRunnable) WhenOverLoad(fn func() error) *DecorateRunnable {
	return OfRunnable(fallback.DecorateRunnableByType[*circuitbreaker.NotPermittedError](runnable.Runnable, fn))
}

func (runnable *DecorateRunnable) WithRetry(entry retry.Retry) *DecorateRunnable {
	return OfRunnable(retry.DecorateRunnable(entry, runnable.Runnable))
}

func (runnable *DecorateRunnable) WhenMaxRetries(fn func() error) *DecorateRunnable {
	return OfRunnable(fallback.DecorateRunnableByType[*retry.MaxRetriesExceeded](runnable.Runnable, fn))
}

func (runnable *DecorateRunnable) WithFallback(
	fn func() error, predicate func(error, any) bool) *DecorateRunnable {
	return OfRunnable(fallback.DecorateRunnable(runnable.Runnable,
		func(ctx fallback.Context[any, any, error]) error { return fn() },
		func(ctx fallback.Context[any, any, error]) (bool, fallback.Context[any, any, error]) {
			return predicate(ctx.Err(), ctx.Panic()), ctx
		}))
}

func (runnable *DecorateRunnable) Decorate() runnable.Runnable {
	return runnable.Runnable
}
