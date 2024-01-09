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

type DecorateConsumer[T any] interface {
	WithBulkhead(bulkhead.Bulkhead) DecorateConsumer[T]
	WhenFull(func(T) error) DecorateConsumer[T]
	WithTimeLimiter(timelimiter.TimeLimiter) DecorateConsumer[T]
	WhenTimeout(func(T) error) DecorateConsumer[T]
	WithRateLimiter(ratelimiter.RateLimiter) DecorateConsumer[T]
	WhenOverRate(func(T) error) DecorateConsumer[T]
	WithCircuitBreaker(circuitbreaker.CircuitBreaker) DecorateConsumer[T]
	WhenOverLoad(func(T) error) DecorateConsumer[T]
	WithRetry(retry.Retry) DecorateConsumer[T]
	WhenMaxRetries(func(T) error) DecorateConsumer[T]
	WithFallback(func(T) error, func(T, error, any) bool) DecorateConsumer[T]
	Decorate() Consumer[T]
}

func OfConsumer[T any](consumer Consumer[T]) DecorateConsumer[T] {
	return &decorateConsumer[T]{consumer}
}

type decorateConsumer[T any] struct {
	Consumer[T]
}

func (fn *decorateConsumer[T]) WithBulkhead(entry bulkhead.Bulkhead) DecorateConsumer[T] {
	return fn.consumer(bulkhead.DecorateConsumer(entry, fn.Consumer))
}

func (fn *decorateConsumer[T]) WhenFull(fallbackFn func(T) error) DecorateConsumer[T] {
	return fn.consumer(fallback.DecorateConsumerByType[T, *bulkhead.FullError](fn.Consumer, fallbackFn))
}

func (fn *decorateConsumer[T]) WithTimeLimiter(entry timelimiter.TimeLimiter) DecorateConsumer[T] {
	return fn.consumer(timelimiter.DecorateConsumer(entry, fn.Consumer))
}

func (fn *decorateConsumer[T]) WhenTimeout(fallbackFn func(T) error) DecorateConsumer[T] {
	return fn.consumer(fallback.DecorateConsumerByType[T, *timelimiter.TimeoutError](fn.Consumer, fallbackFn))
}

func (fn *decorateConsumer[T]) WithRateLimiter(entry ratelimiter.RateLimiter) DecorateConsumer[T] {
	return fn.consumer(ratelimiter.DecorateConsumer(entry, fn.Consumer))
}

func (fn *decorateConsumer[T]) WhenOverRate(fallbackFn func(T) error) DecorateConsumer[T] {
	return fn.consumer(fallback.DecorateConsumerByType[T, *ratelimiter.NotPermittedError](fn.Consumer, fallbackFn))
}

func (fn *decorateConsumer[T]) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) DecorateConsumer[T] {
	return fn.consumer(circuitbreaker.DecorateConsumer(entry, fn.Consumer))
}

func (fn *decorateConsumer[T]) WhenOverLoad(fallbackFn func(T) error) DecorateConsumer[T] {
	return fn.consumer(fallback.DecorateConsumerByType[T, *circuitbreaker.NotPermittedError](fn.Consumer, fallbackFn))
}

func (fn *decorateConsumer[T]) WithRetry(entry retry.Retry) DecorateConsumer[T] {
	return fn.consumer(retry.DecorateConsumer(entry, fn.Consumer))
}

func (fn *decorateConsumer[T]) WhenMaxRetries(fallbackFn func(T) error) DecorateConsumer[T] {
	return fn.consumer(fallback.DecorateConsumerByType[T, *retry.MaxRetriesExceeded](fn.Consumer, fallbackFn))
}

func (fn *decorateConsumer[T]) WithFallback(
	fallbackFn func(T) error, predicate func(T, error, any) bool) DecorateConsumer[T] {
	return fn.consumer(fallback.DecorateConsumer(fn.Consumer,
		func(ctx fallback.Context[T, any, error]) error { return fallbackFn(ctx.Param()) },
		func(ctx fallback.Context[T, any, error]) (bool, fallback.Context[T, any, error]) {
			return predicate(ctx.Param(), ctx.Err(), ctx.Panic()), ctx
		}))
}

func (fn *decorateConsumer[T]) consumer(consumer Consumer[T]) DecorateConsumer[T] {
	fn.Consumer = consumer
	return fn
}

func (fn *decorateConsumer[T]) Decorate() Consumer[T] {
	return fn.Consumer
}
