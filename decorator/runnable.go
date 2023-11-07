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
	return &DecorateRunnable{bulkhead.DecorateRunnable(entry, runnable.Runnable)}
}

func (runnable *DecorateRunnable) WithTimeLimiter(entry timelimiter.TimeLimiter) *DecorateRunnable {
	return &DecorateRunnable{timelimiter.DecorateRunnable(entry, runnable.Runnable)}
}

func (runnable *DecorateRunnable) WithRateLimiter(entry ratelimiter.RateLimiter) *DecorateRunnable {
	return &DecorateRunnable{ratelimiter.DecorateRunnable(entry, runnable.Runnable)}
}

func (runnable *DecorateRunnable) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) *DecorateRunnable {
	return &DecorateRunnable{circuitbreaker.DecorateRunnable(entry, runnable.Runnable)}
}

func (runnable *DecorateRunnable) WithRetry(entry retry.Retry) *DecorateRunnable {
	return &DecorateRunnable{retry.DecorateRunnable(entry, runnable.Runnable)}
}

func (runnable *DecorateRunnable) WithFallback(fn func(error) error) *DecorateRunnable {
	return &DecorateRunnable{fallback.DecorateRunnable(runnable.Runnable, fn)}
}

func (runnable *DecorateRunnable) Decorate() runnable.Runnable {
	return runnable.Runnable
}
