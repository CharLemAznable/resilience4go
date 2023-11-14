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

func (runnable *DecorateRunnable) WhenFull(fn func(*bulkhead.FullError) error) *DecorateRunnable {
	return OfRunnable(fallback.DecorateRunnableDefault(runnable.Runnable, fn))
}

func (runnable *DecorateRunnable) WithTimeLimiter(entry timelimiter.TimeLimiter) *DecorateRunnable {
	return OfRunnable(timelimiter.DecorateRunnable(entry, runnable.Runnable))
}

func (runnable *DecorateRunnable) WhenTimeout(fn func(*timelimiter.TimeoutError) error) *DecorateRunnable {
	return OfRunnable(fallback.DecorateRunnableDefault(runnable.Runnable, fn))
}

func (runnable *DecorateRunnable) WithRateLimiter(entry ratelimiter.RateLimiter) *DecorateRunnable {
	return OfRunnable(ratelimiter.DecorateRunnable(entry, runnable.Runnable))
}

func (runnable *DecorateRunnable) WhenOverRate(fn func(*ratelimiter.NotPermittedError) error) *DecorateRunnable {
	return OfRunnable(fallback.DecorateRunnableDefault(runnable.Runnable, fn))
}

func (runnable *DecorateRunnable) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) *DecorateRunnable {
	return OfRunnable(circuitbreaker.DecorateRunnable(entry, runnable.Runnable))
}

func (runnable *DecorateRunnable) WhenOverLoad(fn func(*circuitbreaker.NotPermittedError) error) *DecorateRunnable {
	return OfRunnable(fallback.DecorateRunnableDefault(runnable.Runnable, fn))
}

func (runnable *DecorateRunnable) WithRetry(entry retry.Retry) *DecorateRunnable {
	return OfRunnable(retry.DecorateRunnable(entry, runnable.Runnable))
}

func (runnable *DecorateRunnable) WhenMaxRetries(fn func(*retry.MaxRetriesExceeded) error) *DecorateRunnable {
	return OfRunnable(fallback.DecorateRunnableDefault(runnable.Runnable, fn))
}

func (runnable *DecorateRunnable) WithFallback(fn func(error) error) *DecorateRunnable {
	return OfRunnable(fallback.DecorateRunnableDefault(runnable.Runnable, fn))
}

func (runnable *DecorateRunnable) Decorate() runnable.Runnable {
	return runnable.Runnable
}
