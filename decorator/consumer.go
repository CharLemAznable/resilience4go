package decorator

import (
	"github.com/CharLemAznable/gofn/consumer"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/fallback"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
)

func OfConsumer[T any](fn func(T) error) *DecorateConsumer[T] {
	return &DecorateConsumer[T]{fn}
}

type DecorateConsumer[T any] struct {
	fn func(T) error
}

func (consumer *DecorateConsumer[T]) WithBulkhead(entry bulkhead.Bulkhead) *DecorateConsumer[T] {
	return consumer.setFn(bulkhead.DecorateConsumer(entry, consumer.fn))
}

func (consumer *DecorateConsumer[T]) WhenFull(fn func(T) error) *DecorateConsumer[T] {
	return consumer.setFn(fallback.DecorateConsumerByType[T, *bulkhead.FullError](consumer.fn, fn))
}

func (consumer *DecorateConsumer[T]) WithTimeLimiter(entry timelimiter.TimeLimiter) *DecorateConsumer[T] {
	return consumer.setFn(timelimiter.DecorateConsumer(entry, consumer.fn))
}

func (consumer *DecorateConsumer[T]) WhenTimeout(fn func(T) error) *DecorateConsumer[T] {
	return consumer.setFn(fallback.DecorateConsumerByType[T, *timelimiter.TimeoutError](consumer.fn, fn))
}

func (consumer *DecorateConsumer[T]) WithRateLimiter(entry ratelimiter.RateLimiter) *DecorateConsumer[T] {
	return consumer.setFn(ratelimiter.DecorateConsumer(entry, consumer.fn))
}

func (consumer *DecorateConsumer[T]) WhenOverRate(fn func(T) error) *DecorateConsumer[T] {
	return consumer.setFn(fallback.DecorateConsumerByType[T, *ratelimiter.NotPermittedError](consumer.fn, fn))
}

func (consumer *DecorateConsumer[T]) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) *DecorateConsumer[T] {
	return consumer.setFn(circuitbreaker.DecorateConsumer(entry, consumer.fn))
}

func (consumer *DecorateConsumer[T]) WhenOverLoad(fn func(T) error) *DecorateConsumer[T] {
	return consumer.setFn(fallback.DecorateConsumerByType[T, *circuitbreaker.NotPermittedError](consumer.fn, fn))
}

func (consumer *DecorateConsumer[T]) WithRetry(entry retry.Retry) *DecorateConsumer[T] {
	return consumer.setFn(retry.DecorateConsumer(entry, consumer.fn))
}

func (consumer *DecorateConsumer[T]) WhenMaxRetries(fn func(T) error) *DecorateConsumer[T] {
	return consumer.setFn(fallback.DecorateConsumerByType[T, *retry.MaxRetriesExceeded](consumer.fn, fn))
}

func (consumer *DecorateConsumer[T]) WithFallback(
	fn func(T) error, predicate func(T, error, any) bool) *DecorateConsumer[T] {
	return consumer.setFn(fallback.DecorateConsumer(consumer.fn,
		func(ctx fallback.Context[T, any, error]) error { return fn(ctx.Param()) },
		func(ctx fallback.Context[T, any, error]) (bool, fallback.Context[T, any, error]) {
			return predicate(ctx.Param(), ctx.Err(), ctx.Panic()), ctx
		}))
}

func (consumer *DecorateConsumer[T]) Decorate() consumer.Consumer[T] {
	return consumer.fn
}

func (consumer *DecorateConsumer[T]) setFn(fn func(T) error) *DecorateConsumer[T] {
	consumer.fn = fn
	return consumer
}
