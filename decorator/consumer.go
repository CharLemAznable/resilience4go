package decorator

import (
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/fallback"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
)

type DecorateConsumer[T any] func(T) error

func OfConsumer[T any](fn func(T) error) DecorateConsumer[T] {
	return fn
}

func (fn DecorateConsumer[T]) WithBulkhead(entry bulkhead.Bulkhead) DecorateConsumer[T] {
	return bulkhead.DecorateConsumer(entry, fn)
}

func (fn DecorateConsumer[T]) WhenFull(fallbackFn func(T) error) DecorateConsumer[T] {
	return fallback.DecorateConsumerByType[T, *bulkhead.FullError](fn, fallbackFn)
}

func (fn DecorateConsumer[T]) WithTimeLimiter(entry timelimiter.TimeLimiter) DecorateConsumer[T] {
	return timelimiter.DecorateConsumer(entry, fn)
}

func (fn DecorateConsumer[T]) WhenTimeout(fallbackFn func(T) error) DecorateConsumer[T] {
	return fallback.DecorateConsumerByType[T, *timelimiter.TimeoutError](fn, fallbackFn)
}

func (fn DecorateConsumer[T]) WithRateLimiter(entry ratelimiter.RateLimiter) DecorateConsumer[T] {
	return ratelimiter.DecorateConsumer(entry, fn)
}

func (fn DecorateConsumer[T]) WhenOverRate(fallbackFn func(T) error) DecorateConsumer[T] {
	return fallback.DecorateConsumerByType[T, *ratelimiter.NotPermittedError](fn, fallbackFn)
}

func (fn DecorateConsumer[T]) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) DecorateConsumer[T] {
	return circuitbreaker.DecorateConsumer(entry, fn)
}

func (fn DecorateConsumer[T]) WhenOverLoad(fallbackFn func(T) error) DecorateConsumer[T] {
	return fallback.DecorateConsumerByType[T, *circuitbreaker.NotPermittedError](fn, fallbackFn)
}

func (fn DecorateConsumer[T]) WithRetry(entry retry.Retry) DecorateConsumer[T] {
	return retry.DecorateConsumer(entry, fn)
}

func (fn DecorateConsumer[T]) WhenMaxRetries(fallbackFn func(T) error) DecorateConsumer[T] {
	return fallback.DecorateConsumerByType[T, *retry.MaxRetriesExceeded](fn, fallbackFn)
}

func (fn DecorateConsumer[T]) WithFallback(
	fallbackFn func(T) error, predicate func(T, error, any) bool) DecorateConsumer[T] {
	return fallback.DecorateConsumer(fn,
		func(ctx fallback.Context[T, any, error]) error { return fallbackFn(ctx.Param()) },
		func(ctx fallback.Context[T, any, error]) (bool, fallback.Context[T, any, error]) {
			return predicate(ctx.Param(), ctx.Err(), ctx.Panic()), ctx
		})
}

func (fn DecorateConsumer[T]) Decorate() func(T) error {
	return fn
}
