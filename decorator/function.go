package decorator

import (
	"github.com/CharLemAznable/gofn/function"
	"github.com/CharLemAznable/resilience4go/bulkhead"
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

func (function *DecorateFunction[T, R]) WithTimeLimiter(entry timelimiter.TimeLimiter) *DecorateFunction[T, R] {
	return OfFunction(timelimiter.DecorateFunction(entry, function.Function))
}

func (function *DecorateFunction[T, R]) WithRateLimiter(entry ratelimiter.RateLimiter) *DecorateFunction[T, R] {
	return OfFunction(ratelimiter.DecorateFunction(entry, function.Function))
}

func (function *DecorateFunction[T, R]) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) *DecorateFunction[T, R] {
	return OfFunction(circuitbreaker.DecorateFunction(entry, function.Function))
}

func (function *DecorateFunction[T, R]) WithRetry(entry retry.Retry) *DecorateFunction[T, R] {
	return OfFunction(retry.DecorateFunction(entry, function.Function))
}

func (function *DecorateFunction[T, R]) WithFallback(fn func(error) (R, error)) *DecorateFunction[T, R] {
	return OfFunction(fallback.DecorateFunction(function.Function, fn))
}

func (function *DecorateFunction[T, R]) Decorate() function.Function[T, R] {
	return function.Function
}
